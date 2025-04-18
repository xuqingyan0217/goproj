package msgTransfer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"

	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/apps/task/mq/mq"
	"easy-chat/pkg/bitmap"
	"easy-chat/pkg/constants"
	"sync"
	"time"
)

// 设置相关配置的默认值
var (
	GroupMsgReadRecordDelayTime = time.Second
	GroupMsgReadRecordDelayCount = 10
)

const (
	GroupMsgReadHandlerAtTransfer = iota // 定义为0，表示不开启
	GroupMsgReadHandlerDelayTransfer  // 定义延迟
)

type MsgReadTransfer struct {
	*baseMsgTransfer

	cache.Cache
	mu sync.Mutex
	// 使用map存储是考虑可能会有其它的群
	groupMsgs map[string]*groupMsgRead
	push chan *ws.Push
}

// NewMsgReadTransfer 创建并返回一个新的MsgReadTransfer实例，该实例实现了kq.ConsumeHandler接口。
// 它负责处理消息消费的逻辑，基于提供的服务上下文（svc）来初始化。
func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	// 返回一个新的MsgReadTransfer实例，该实例包装了NewBaseMsgTransfer创建的基础消息传输逻辑。
	// 通过这种方式，MsgReadTransfer获得了基础消息传输的功能，同时可以添加或覆盖具体的消息读取处理逻辑。
	m := &MsgReadTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
		groupMsgs: make(map[string]*groupMsgRead, 1),
		push: make(chan *ws.Push, 1),
	}
	// 不等于0，表示开启延迟处理
	if svc.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount > 0 {
			GroupMsgReadRecordDelayCount = svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount
		}
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			GroupMsgReadRecordDelayTime = time.Duration(svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime)*time.Second
		}
	}
	go m.transfer()

	return m
}

// Consume 消息消费处理函数，处理消息已读请求并返回错误。
func (m *MsgReadTransfer) Consume(key, value string) error {
	m.Info("MsgReadTransfer", value)

	var (
		data mq.MsgMarkRead
		ctx  = context.Background()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 业务处理 -- 更新

	// map[消息id]已读记录，key-value都选择string类型，因为消息需要在mq再传到ws，避免问题，选择string
	// 我们选择该map作为更新聊天记录的返回值，这样用户就可以通过该map知道该条甚至多条消息已读情况
	// todo :业务逻辑
	readRecords,err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ReadRecords: 	readRecords,
		ContentType:    constants.ContentMarkRead,
		//RecvIds:        data.RecvIds,
		//SendTime:       data.SendTime,
		//MType:          data.MType,
		//Content:        data.Content,
	}
	switch push.ChatType {
	case constants.SingleChatType:
		// 直接推送
		// 先把消息放到推送管道中，然后把推送消息换成异步处理，由下面这个管道进行通信
		m.push <- push
	case constants.GroupChatType:
		// 判断是否开启合并消息的处理，若没有，直接给到channel里
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			m.push <- push
		}
		m.mu.Lock()
		defer m.mu.Unlock()
		// 因为合并，发送方id置为空就行
		push.SendId = ""
		// 判断是否已经存在该消息的记录，若存在，则合并；没有则创建
		if _,ok := m.groupMsgs[push.ConversationId];ok {
			m.Infof("mergePush", push.ConversationId)
			m.groupMsgs[push.ConversationId].mergePush(push)
		} else {
			m.Infof("NewGroupMsgRead", push.ConversationId)
			m.groupMsgs[push.ConversationId] = NewGroupMsgRead(push, m.push)
		}
	}
	// 处理完请求之后，调用基础消息传输的Transfer方法，将消息推送给客户端。
	return nil
}

// UpdateChatLogRead 更新聊天记录已读状态；针对websocket推送的消息做处理
// 该函数接收一个消息标记已读的MQ消息，更新相应的聊天记录，并返回更新后的已读状态
func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
    // 初始化返回的已读状态映射
    res := make(map[string]string)
    // 下面我们需要根据id集合来获取对应的聊天记录(消息)
    chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
    if err != nil {
        return nil, err
    }
    // 处理已读
    for _, chatLog := range chatLogs {
        // 获取已读记录，并判断是否已读，如果已读则跳过
        if chatLog.ReadRecords != nil && bitmap.Load(chatLog.ReadRecords).IsSet(data.SendId) {
            continue
        }
        // 更新已读记录
        // 依据私聊和群聊分别处理
        switch chatLog.ChatType {
        case constants.SingleChatType:
			// 私聊只有一个字节就可以了，我们在这里赋值为1，表示请求方对私聊的消息已读
            chatLog.ReadRecords = []byte{1}
        case constants.GroupChatType:
			// 群聊则需要将发送请求者的设置到bitmap上
            readRecords := bitmap.Load(chatLog.ReadRecords)
            readRecords.Set(data.SendId)
			// 将更新后的已读记录导出为字节数组存到chatLog里
            chatLog.ReadRecords = readRecords.Export()
        }
        // 将更新后的已读记录编码后存入结果映射
        res[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)

        // 更新数据库中的聊天记录已读状态
        err = m.svcCtx.ChatLogModel.UpdateMarkRead(ctx, chatLog.ID, chatLog.ReadRecords)
        if err != nil {
            return nil, err
        }
    }
    // 返回更新后的已读状态映射
    return res, nil
}

// 异步推送消息
// MsgReadTransfer 的 transfer 方法用于处理消息传输逻辑。
// 该方法不断监听并处理推送过来的消息，该消息可以是群聊可以是私聊，根据不同的聊天类型和配置决定是否进行消息传输和如何处理消息。
func (m *MsgReadTransfer) transfer() {
    // 从推送通道中循环读取消息
    for push := range m.push {
		fmt.Println("00100000000000000000000")
        // 检查消息接收者ID，如果非空则执行传输操作
        if push.RecvId != "" || len(push.RecvIds) > 0 {
			fmt.Println("00200000000000000000000")
            // 执行消息传输，如果发生错误则记录错误信息
            if err := m.Transfer(context.Background(), push); err != nil {
                m.Errorf("m transfer err %v push %v", err, push)
            }
        }
		fmt.Println("003000000000000000000")
		// 对于单聊类型的消息，直接进入下一轮循环，所以说，其实私聊就走到这里就ok了，再往后是留给捕获群聊消息的
        if push.ChatType == constants.SingleChatType {
            continue
        }

        // 该语句往下都是群聊类型才能到的
		// 下面这一句很简单，它主要针对群聊消息下，因关闭缓存的原因，所以群聊消息也是无需缓存，直接推送的，那么在下面的删除操作就不需要了，所以要跳
		// 看似是和群聊的第一个判断重复了，但其实就是为了在这种情况下跳出流程，防止执行下面不该执行的东西。
        if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
            continue
        }
        // 对于群聊类型的消息，在数据推送之后，根据特定条件清空消息数据
		// 开启锁，此时该协程就会因锁而阻塞了
        m.mu.Lock()
        if _, ok := m.groupMsgs[push.ConversationId]; ok && m.groupMsgs[push.ConversationId].IsIdle() {
            m.groupMsgs[push.ConversationId].clear()
            delete(m.groupMsgs, push.ConversationId)
        }
		m.mu.Unlock()
    }
}