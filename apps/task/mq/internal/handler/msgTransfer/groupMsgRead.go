package msgTransfer

import (
	"github.com/zeromicro/go-zero/core/logx"

	"easy-chat/apps/im/ws/ws"
	"easy-chat/pkg/constants"
	"sync"
	"time"
)

// 仅代表一个群里面的群聊消息，不能是多个群
type groupMsgRead struct {
	mu sync.Mutex
	ConversationId string
	// 记录消息
	push *ws.Push
	// 推送消息,连接着groupMsgRead和msgReadTransfer的推送
	pushCh chan *ws.Push
	// 统计数量
	count int
	// 上次推送时间
	pushTime time.Time

	done chan struct{}
}

func NewGroupMsgRead(push *ws.Push, pushCh chan *ws.Push) *groupMsgRead {
	m := &groupMsgRead{
		ConversationId: push.ConversationId,
		push:     push,
		pushCh:   pushCh,
		count:    1,
		pushTime: time.Now(),
		done:     make(chan struct{}),
	}
	go m.transfer()
	return m
}

//消息合并
func (m *groupMsgRead) mergePush(push *ws.Push) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.count++
	// 遍历形参，将其内容添加到该对象字节的push里面
	for msgId, read := range push.ReadRecords {
		m.push.ReadRecords[msgId] = read
	}
}

// 合并后的检测
// transfer 方法用于处理群消息的读取状态传输逻辑。
// 它主要负责在满足一定条件时，将消息的读取状态发送给服务器。
// 该方法在内部会根据读取状态的变化频率和数量来决定何时发送状态更新。
func (m *groupMsgRead) transfer() {
    // 1. 超时则发送
    // 2. 超量则发送

    // 初始化一个定时器，用于跟踪读取状态的超时时间
    timer := time.NewTimer(GroupMsgReadRecordDelayTime / 2)
    defer timer.Stop()

    for {
        select {
        case <-m.done:
            return
        case <-timer.C:
            m.mu.Lock()

            // 检查是否到推送时间，计算距离上次推送的时间间隔
            pushTime := m.pushTime
            val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)
            push := m.push
            // 在符合检测条件内，无需推送
            if val > 0 && m.count < GroupMsgReadRecordDelayCount || push == nil {
                if val > 0 {
                    // 重置定时器
                    timer.Reset(val)
                }
                m.mu.Unlock()
                continue
            }
            // 到这说明超出检测条件了，需要推送了
            // 获取推送的新的时间点
            m.pushTime = time.Now()
            // 推送前把数据清空
            m.push = nil
            m.count = 0
            // 重置定时器
            timer.Reset(GroupMsgReadRecordDelayTime / 2)
            m.mu.Unlock()

            // 日志记录：超过合并的条件，推送消息
            logx.Infof("超过合并的条件，推送 %v", push)
            m.pushCh <- push
        default:
            m.mu.Lock()
            // 检查消息数量是否达到推送条件
            if m.count >= GroupMsgReadRecordDelayCount {
                push := m.push
                m.pushTime = time.Now()
                // 推送前把数据清空
                m.push = nil
                m.count = 0
                m.mu.Unlock()

                // 日志记录：消息数量超过合并条件，推送消息
                logx.Infof("default 超过合并的条件，推送 %v", push)
                m.pushCh <- push
                continue
            }
            // 判断是否为空闲状态
            if m.isIdle() {
                m.mu.Unlock()
                // 发送空闲信号，以释放msgReadTransfer
                m.pushCh <- &ws.Push{
                    ConversationId: m.ConversationId,
                    ChatType:       constants.GroupChatType,
                }
                continue
            }
            // 不为空闲状态
            m.mu.Unlock()
            // 短暂休眠，避免频繁操作
            tempDelay := GroupMsgReadRecordDelayTime / 4
            if tempDelay > time.Second {
                tempDelay = time.Second
            }
            time.Sleep(tempDelay)
        }
    }
}


// 判断活跃状态
func (m *groupMsgRead) isIdle() bool {
	// 获取上一次推送时间
	pushTime := m.pushTime
	// 计算时间差，看是否超时，但这里*2，是为了给点容错
	val := GroupMsgReadRecordDelayTime * 2 - time.Since(pushTime)
	// 过期了，空闲
	if val <= 0 && m.push == nil && m.count == 0 {
		return true
	}
	return false
}

func (m *groupMsgRead) IsIdle() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isIdle()
}

// 清空数据
func (m *groupMsgRead) clear()  {
	select {
	case <-m.done:
	default:
		close(m.done)
	}

	m.push = nil
}