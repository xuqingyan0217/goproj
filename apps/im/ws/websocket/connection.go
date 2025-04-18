/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
	"time"
)

type Conn struct {
	idleMu sync.Mutex

	Uid string

	*websocket.Conn
	s *Server

	idle              time.Time
	maxConnectionIdle time.Duration

	messageMu      sync.Mutex
	readMessage    []*Message
	readMessageSeq map[string]*Message

	message chan *Message

	done chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {

	var responseHeader http.Header
	if protocol := r.Header.Get("Sec-Websocket-Protocol"); protocol != "" {
		responseHeader = http.Header{
			"Sec-Websocket-Protocol": []string{protocol},
		}
	}

	c, err := s.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		s.Errorf("upgrade err %v", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,
		readMessage:       make([]*Message, 0, 2),
		readMessageSeq:    make(map[string]*Message, 2),
		message:           make(chan *Message, 1),
		done:              make(chan struct{}),
	}
	go conn.keepalive()
	return conn
}

func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()
	// 读队列中的某一特定id项是否有值，有的话说明之前已经做过一次记录了，接下来要执行更新。
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		if len(c.readMessage) == 0 {
			// 当前判断是：可能刚才查看id序号是存在的，但是此时消息数组里面已经没有值了，所以说我们需要有这一步判断，这种时候直接返回
			return
		}

		// msg.AckSeq > m.AckSeq，判断当前id序号是否符合要求
		// m是我们之前记录过的消息；而msg是新来的消息，我们要进行更新，但在更新前，先验证一下。
		if m.AckSeq >= msg.AckSeq {
			// 没有进行ack的确认, 重复
			return
		}
		// 到这就没错了，执行更新消息吧，更新序号
		c.readMessageSeq[msg.Id] = msg
		return
	}
	// 到这说明还没有进行消息的记录，准备记录
	// 避免客户端重复发送多余的ack消息，
	if msg.FrameType == FrameAck {
		return
	}
	// 更新
	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg

}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()
	if err != nil {
		logx.Info("ReadMessage error: %v", err)
		return
	}

	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Now()
	logx.Info("ReadMessage: idle time reset to zero, conn uid:",c.Uid)
	return
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()

	// 方法是并不安全
	err := c.Conn.WriteMessage(messageType, data)
	if err != nil {
		logx.Info("WriteMessage error: %v", err)
		return err
	}

	c.idle = time.Now()
	logx.Infof("WriteMessage: idle time updated to %v, conn uid:%v", c.idle, c.Uid)
	return nil
}


func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	return c.Conn.Close()
}

func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		if c.Uid == "root" {
			return
		}
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() { // The connection is non-idle.
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				// The connection has been idle for a duration of keepalive.MaxConnectionIdle or more.
				// Gracefully close the connection.
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			return
		}
	}
}


