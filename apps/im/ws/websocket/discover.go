package websocket

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"
)

// 服务发现机制【该方式是去中心化，自己在内部实现服务发现整套机制】
// 该机制主要针对用户而设立
//
// 用户连接后，会将用户信息与服务器ip一起绑定注册到某一个位置比如redis
// 当用户发送信息的时候，根据发送的目标从记录位置中获取绑定关系，查找相应服务并发送

// Discover 是一个接口，定义了服务发现和用户绑定相关的方法
type Discover interface {
	// Register 用于注册服务到服务发现系统
	// 参数 serverAddr 是服务的地址
	// 返回 error 如果注册过程中发生错误
	Register(serverAddr string) error

	// BoundUser 用于将用户绑定到特定的服务
	// 参数 uid 是用户的唯一标识符
	// 返回 error 如果绑定过程中发生错误
	BoundUser(uid string) error

	// RelieveUser 用于解除用户与服务的绑定关系
	// 参数 uid 是用户的唯一标识符
	// 返回 error 如果解除绑定过程中发生错误
	RelieveUser(uid string) error

	// Transpond 用于转发消息到指定用户或所有用户
	// 参数 msg 是要转发的消息，可以是任意类型
	// 参数 uid 是一个可变参数，指定消息接收者的用户标识符，如果为空，则消息广播给所有用户
	// 返回 error 如果转发过程中发生错误
	Transpond(msg interface{}, uid ...string) error
}

// 针对如上接口，定义一个空的实现
type nopDiscover struct {
	serverAddr string
}

// Register 注册服务
func (d *nopDiscover) Register(serverAddr string) error { return nil }
// BoundUser 绑定用户
func (d *nopDiscover) BoundUser(uid string) error { return nil }
// RelieveUser 解除绑定
func (d *nopDiscover) RelieveUser(uid string) error { return nil }
// Transpond 转发消息
func (d *nopDiscover) Transpond(msg interface{}, uid ...string) error { return nil }

// 定义一个非空的实现
type redisDiscover struct {
	serverAddr string
	auth http.Header
	srvKey string
	boundUserKey string
	redis *redis.Redis
	clients map[string]Client
}
// NewRedisDiscover 创建一个新的redisDiscover实例，用于服务发现和客户端管理。
// 参数:
//   auth http.Header: 认证信息，用于客户端请求。
//   srvKey string: 服务的唯一键，用于在Redis中标识服务。
//   redisCfg redis.RedisConf: Redis配置信息，用于连接Redis服务器。
// 返回值:
//   *redisDiscover: redisDiscover的实例。
func NewRedisDiscover(auth http.Header, srvKey string, redisCfg redis.RedisConf) *redisDiscover {
	return &redisDiscover{
		srvKey: srvKey,
		boundUserKey: fmt.Sprintf("%s.%s", srvKey, "boundUserKey"),
		redis: redis.MustNewRedis(redisCfg),
		clients: make(map[string]Client),
		auth: auth,
	}
}

// Register 用于注册服务地址到Redis。
// 参数:
//   serverAddr string: 服务的地址。
// 返回值:
//   (err error): 错误信息，如果有的话。
func (d *redisDiscover) Register(serverAddr string) (err error) {
	d.serverAddr = serverAddr
	// 服务列表：redis存储用set
	go d.redis.Set(d.srvKey, serverAddr)
	return
}

// BoundUser 用于绑定用户到当前服务。
// 参数:
//   uid string: 用户的唯一标识符。
// 返回值:
//   (err error): 错误信息，如果有的话。
func (d *redisDiscover) BoundUser(uid string) (err error) {
	// 用户绑定
	exists, err := d.redis.Hexists(d.boundUserKey, uid)
	if err != nil {
		return err
	}
	if exists {
		// 存在绑定关系
		return nil
	}
	// 绑定
	return d.redis.Hset(d.boundUserKey, uid, d.serverAddr)
}

// RelieveUser 用于解除用户与服务的绑定关系。
// 参数:
//   uid string: 用户的唯一标识符。
// 返回值:
//   (err error): 错误信息，如果有的话。
func (d *redisDiscover) RelieveUser(uid string) (err error) {
	_, err = d.redis.Hdel(d.boundUserKey, uid)
	return
}

// Transpond 用于转发消息到绑定的用户。
// 参数:
//   msg interface{}: 要发送的消息。
//   uids ...string: 用户的唯一标识符列表。
// 返回值:
//   (err error): 错误信息，如果有的话。
func (d *redisDiscover) Transpond(msg interface{}, uids ...string) (err error) {
	// 遍历用户ID列表，以处理每个用户ID
	for _, uid := range uids {
	    // 从Redis中获取与用户ID关联的服务地址
	    srvAddr, err := d.redis.Hget(d.boundUserKey, uid)
	    if err != nil {
	        // 如果获取服务地址时发生错误，则返回错误
	        return err
	    }

	    // 尝试获取与服务地址关联的客户端
	    srvClient, ok := d.clients[srvAddr]
	    if !ok {
	        // 如果不存在与服务地址关联的客户端，则创建一个新的客户端
	        srvClient = d.createClient(srvAddr)
	    }

	    // 打印日志信息，指示正在处理的用户ID和服务地址
	    fmt.Println("redis transpand -》 ", srvAddr, " uid ", uid)

	    // 向服务客户端发送消息
	    if err := d.send(srvClient, msg, uid); err != nil {
	        // 如果发送消息时发生错误，则返回错误
	        return err
	    }
	}
	return
}

// send 发送消息到指定的客户端。
// 参数:
//   srvClient Client: 服务客户端实例。
//   msg interface{}: 要发送的消息。
//   uid string: 用户的唯一标识符。
// 返回值:
//   error: 错误信息，如果有的话。
func (d *redisDiscover) send(srvClient Client, msg interface{}, uid string) error {
	return srvClient.Send(Message{
		FrameType: FrameTranspond,
		TranspondUid: uid,
		Data: msg,
	})
}

// createClient 创建一个新的客户端实例。
// 参数:
//   srvAddr string: 服务的地址。
// 返回值:
//   Client: 新创建的客户端实例。
func (d *redisDiscover) createClient(srvAddr string) Client {
	return NewClient(srvAddr, WithClientHeader(d.auth))
}
