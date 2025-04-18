/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/threading"
	"time"

	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type AckType int

const (
	NoAck AckType = iota
	OnlyAck
	RigorAck
)

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	}

	return "NoAck"
}

type Server struct {
	sync.RWMutex

	*threading.TaskRunner

	opt            *serverOption
	authentication Authentication

	routes map[string]HandlerFunc
	addr   string
	patten string

	connToUser map[*Conn]string
	userToConn map[string]*Conn

	upgrader websocket.Upgrader
	logx.Logger

	listenOn string
	discover Discover
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	s := &Server{
		routes:   make(map[string]HandlerFunc),
		addr:     addr,
		patten:   opt.patten,
		opt:      &opt,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},

		authentication: opt.Authentication,

		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),

		Logger: logx.WithContext(context.Background()),

		TaskRunner: threading.NewTaskRunner(opt.concurrency),

		listenOn: FigureOutListenOn(addr),
		discover: opt.Discover,
	}
	// 存在服务发现，采用分布式im通信的时候; 默认不做任何处理
	// 这里if判断仅仅是个提示，其实就那一句就可以了
	/*if s.discover == nil {
		fmt.Println("不使用服务发现")
		s.discover.Register(fmt.Sprintf("%s", s.listenOn))
	} else {
		fmt.Println("使用服务发现")
		s.discover.Register(fmt.Sprintf("%s", s.listenOn))
	}*/
	s.discover.Register(fmt.Sprintf("%s", s.listenOn))

	return s
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != any(nil) {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}
	//conn, err := s.upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	s.Errorf("upgrade err %v", err)
	//	return
	//}

	if !s.authentication.Auth(w, r) {
		//conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint("不具备访问权限")))
		s.Send(&Message{FrameType: FrameData, Data: fmt.Sprint("不具备访问权限")}, conn)
		conn.Close()
		return
	}

	// 记录连接
	s.addConn(conn, r)

	// 处理连接
	go s.handlerConn(conn)
}

// 根据连接对象执行任务处理
func (s *Server) handlerConn(conn *Conn) {

	uids := s.GetUsers(conn)
	conn.Uid = uids[0]
	logx.Info("用户登入---uid:", conn.Uid)

	// 如果存在服务发现则进行注册；默认不做任何处理
	s.discover.BoundUser(conn.Uid)

	// 处理任务
	go s.handlerWrite(conn)

	if s.isAck(nil) {
		go s.readAck(conn)
	}

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			s.Close(conn)
			return
		}
		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v, msg %v", err, string(msg))
			//s.Close(conn)
			//return
			continue
		}

		// 依据消息进行处理
		if s.isAck(&message) {
			s.Infof("conn message read ack msg %v", message)
			conn.appendMsgMq(&message)
		} else {
			conn.message <- &message
		}
	}
}

func (s *Server) isAck(message *Message) bool {
	if message == nil {
		return s.opt.ack != NoAck
	}
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck && message.FrameType != FrameTranspond
}

// 读取消息的ack
func (s *Server) readAck(conn *Conn) {
	// 失败次数,发送失败超过一定次数后，
	/*send := func(msg *Message, conn *Conn) error {
		err := s.Send(msg, conn)
		if err == nil {
			// 一般来说，服务端都能发送出去，直接返回了，
			return nil
		}

		s.Errorf("message ack OnlyAck send err %v message %v", err, msg)
		conn.readMessage[0].errCount++
		conn.messageMu.Unlock()

		tempDelay := time.Duration(200*conn.readMessage[0].errCount) * time.Microsecond
		if max := 1 * time.Second; tempDelay > max {
			tempDelay = max
		}
		time.Sleep(tempDelay)
		return err
	}*/

	for {
		select {
		case <-conn.done:
			s.Infof("close message ack uid %v ", conn.Uid)
			return
		default:
		}
		// 到这说明连接没有关闭

		// 从队列中读取新的消息
		conn.messageMu.Lock()

		// 如果队列中没有值，那就解锁，并直接结束本次循环
		if len(conn.readMessage) == 0 {
			conn.messageMu.Unlock()
			// 增加睡眠，为了让任务更好切换（100ms）
			time.Sleep(100 * time.Microsecond)
			continue
		}
		// 如果队列中有值
		// 读取第一条
		message := conn.readMessage[0]

		// 判断ack的方式，两种，一次应答和两次应答
		switch s.opt.ack {
		case OnlyAck:
			// 直接给客户端回复，发送ack确认消息，结束本次应答
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq + 1,
			}, conn)
			// 进行业务处理
			// 把消息从队列中移除，选择切片的方式，直接从第二个开始作为新切片，这样往后的都会变成第一个
			conn.readMessage = conn.readMessage[1:]
			// 释放锁
			conn.messageMu.Unlock()
			// 将消息传入到管道中供另一个协程处理
			conn.message <- message
		case RigorAck: // 严格模式
			// 先回
			if message.AckSeq == 0 {
				// 还未确认
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].ackTime = time.Now()
				s.Send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn)
				s.Infof("message ack RigorAck send mid %v, seq %v , time%v", message.Id, message.AckSeq,
					message.ackTime)
				conn.messageMu.Unlock()
				continue
			}

			// 再验证

			// 1. 客户端返回结果，再一次确认
			// 得到客户端的序号
			// 此处为什么可以从map中获取到最新序列呢，这是因为客户端的应答消息在appendMsgMq方法那里实现了map更新，所以说数组是旧的，map是新的
			msgSeq := conn.readMessageSeq[message.Id]
			if msgSeq.AckSeq > message.AckSeq {
				// 确认
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				s.Infof("message ack RigorAck success mid %v", message.Id)
				continue
			}

			// 2. 客户端没有确认，考虑是否超过了ack的确认时间
			val := s.opt.ackTimeout - time.Since(message.ackTime)
			if !message.ackTime.IsZero() && val <= 0 {
				//		2.2 超过结束确认
				delete(conn.readMessageSeq, message.Id)
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				continue
			}
			//		2.1 未超过，重新发送
			conn.messageMu.Unlock()
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq,
			}, conn)
			// 睡眠一定的时间
			time.Sleep(3 * time.Second)
		}
	}
}

// 任务的处理
func (s *Server) handlerWrite(conn *Conn) {
	for {
		select {
		case <-conn.done:
			// 连接关闭
			return
		case message := <-conn.message:
			switch message.FrameType {
			case FramePing:
				s.Send(&Message{FrameType: FramePing}, conn)
			case FrameData:
				// 根据请求的method分发路由并执行
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				} else {
					s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)}, conn)
					//conn.WriteMessage(&Message{}, []byte(fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)))
				}
			}

			if s.isAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}
}

func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 验证用户是否之前登入过
	if c := s.userToConn[uid]; c != nil {
		// 关闭之前的连接
		c.Close()
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	/*if _, ok := s.userToConn[uid]; !ok {

	}*/
	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) ([]*Conn, []string) {
	if len(uids) == 0 {
		return nil, nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	noExistUids := make([]string, 0, len(uids))
	for _, uid := range uids {
		// 新增判断，判断当前map里面是否有值，有值说明用户发送的对象在当前服务器上，没有则说明用户离线或者不在一台机器上
		if _, ok := s.userToConn[uid]; !ok {
			noExistUids = append(noExistUids, uid)
		}
		res = append(res, s.userToConn[uid])
	}
	return res, noExistUids
}

func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经被关闭
		return
	}

	fmt.Printf("关闭与%s的连接,连接对象:%v\n", uid, conn)

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	s.discover.RelieveUser(uid)

	conn.Close()
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}
	// 对GetConns方法进行简单改变
	conns, noExistUids := s.GetConns(sendIds...)
	fmt.Println("发送消息 noExistUids : ", noExistUids, " sendIds ", sendIds)
	// 发送当前服务的
	err := s.Send(msg, conns...)
	if err != nil {
		return err
	}
	// 转发由其他服务处理
	return s.discover.Transpond(msg, noExistUids...)
}


func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("停止服务")
}
