package interceptor

import (
	"context"
	"easy-chat/pkg/xerr"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type Idempotent interface {
	// Identify 获取请求的标识，该标识是将TKey组合成DKey
	Identify(ctx context.Context, method string) string
	// IsIdempotentMethod 是否支持幂等性，参数就是作为map的key，该map是结构体对象里面的
	IsIdempotentMethod(fullMethod string) bool
	// TryAcquire 幂等性的验证，就是去访问redis，第一次访问会上锁，后续重试因锁而无法访问【过期时间足够大】，如此实现幂等
	TryAcquire(ctx context.Context, id string) (resp interface{}, isAcquire bool)
	// SaveResp 执行之后结果的保存，做缓存。
	SaveResp(ctx context.Context, id string, resp interface{}, respErr error) error
}

var (
	// TKey 请求任务标识，作为刚生成的key，有context生成，该id每次都不一样
	TKey = "easy-chat-idempotence-task-id"
	// DKey 设置rpc调度中rpc请求的标识，上面的TKey和服务method的组合形成该key，也是最终的key
	DKey = "easy-chat-idempotence-dispatch-key"
)

// ContextWithVal 上下文设置的最初请求唯一id，该方法在中间件里面调用
func ContextWithVal(ctx context.Context) context.Context {
	// 设置请求的id，唯一性id主要就是来自于该方法
	return context.WithValue(ctx, TKey, utils.NewUuid())
}

// NewIdempotenceClient 客户端拦截器
func NewIdempotenceClient(idempotent Idempotent) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 获取唯一的key，此时的key就是方法封装过了的最终id值
		identify := idempotent.Identify(ctx, method)

		// 通过父ctx的方法，将最终的id值设置到子ctx中
		ctx = metadata.NewOutgoingContext(ctx, map[string][]string{
			DKey: {identify},
		})
		// 执行rpc请求
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// NewIdempotenceServer 服务端拦截器
func NewIdempotenceServer(idempotent Idempotent) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 获取请求的id，注意，服务端就要从ctx里面获取了
		identify := metadata.ValueFromIncomingContext(ctx, DKey)
		// 如下是不进行幂等性处理
		// 其中第一项判断是传入的参数是map的key，这个map在结构体对象里面，在初始化里面往map里面添加了一个social服务的key，值为true
		// 而这个info.FullMethod就能够识别到我们初始化里面写的key类型，它也可以识别像user等的服务；
		// 所以说，传入它就能进行判断相关服务是否进行幂等，都在map我们手动添加，所以，只要在map里面有，就能由info.FullMethod判断幂等
		if !idempotent.IsIdempotentMethod(info.FullMethod) || len(identify) == 0 {
			// 要么是没有该标识性id，要么是方法不支持幂等
			return handler(ctx, req)
		}

		// 进行幂等性
		fmt.Println("----" ,"请求进入 幂等性处理 ", identify)
		// 判断当前请求的处理状态，是已处理还是正在处理，这里面就是幂等的核心，redis锁
		r, isAcquire := idempotent.TryAcquire(ctx, identify[0])
		if isAcquire {
			resp, err = handler(ctx, req)
			fmt.Println("----执行任务", identify)

			if err := idempotent.SaveResp(ctx, identify[0], resp, err); err != nil {
				return resp, err
			}
			return resp, nil
		}
		// 到这说明不走if，任务已经在执行了
		fmt.Println("----任务正在执行", identify)

		if r != nil {
			fmt.Println("----任务已经执行完了", identify)
			return r, nil
		}
		// 否则，任务可能还在执行
		return nil, errors.WithStack(xerr.New(int(codes.DeadlineExceeded), fmt.Sprintf("任务还在执行，id %v",identify)))
	}
}

// 定义一个默认的初始化结构体对象去实现接口，并使用它定义一个默认的客户端拦截器初始化
var (
	DefaultIdempotent = new(defaultIdempotent)
	DefaultIdempotentClient = NewIdempotenceClient(DefaultIdempotent)
)
// 定义结构体对象
type defaultIdempotent struct {
	// 获取和设置请求的id
	*redis.Redis
	// 注意存储
	*collection.Cache
	// 设置方法对幂等的支持
	method map[string]bool
}
// NewDefaultIdempotent 初始化结构体对象
func NewDefaultIdempotent(c redis.RedisConf) Idempotent {
	cache, err := collection.NewCache(60 * 60)
	if err != nil {
		panic(any(err))
	}
	return &defaultIdempotent{
		Redis: redis.MustNewRedis(c),
		Cache: cache,
		method: map[string]bool{
			// 可依据业务往里添加user服务，im服务等，格式一样就行，这里是0是1直接决定是否开启幂等
			"/social.social/GroupCreate": true,
		},
	}
}

// 实现接口方法

// Identify 生成请求的唯一标识
// 该方法结合了上下文中的标识和方法名，用于在RPC调用中唯一地标识一个请求
func (d *defaultIdempotent) Identify(ctx context.Context, method string) string {
	// 该key是最初父ctx设置到子ctx里面的
    id := ctx.Value(TKey)
    // 生成最后的请求id
    rpcId := fmt.Sprintf("%v.%s",id, method)
    return rpcId
}

// IsIdempotentMethod 判断方法是否是幂等的
// 该方法检查给定的完整方法名是否在幂等方法的映射中
func (d *defaultIdempotent) IsIdempotentMethod(fullMethod string) bool {
    return d.method[fullMethod]
}

// TryAcquire 尝试获取锁并处理请求
// 该方法首先尝试在Redis中设置一个键（如果尚未存在），以确保请求的幂等性如果键已经存在，则获取并返回缓存的响应
func (d *defaultIdempotent) TryAcquire(ctx context.Context, id string) (resp interface{}, isAcquire bool) {
    // 基于redis锁设置
    // id: 锁的唯一标识符
    // "1": 锁的值，通常用于标识锁的状态
    // 60 * 60: 锁的有效期，这里设置为60分钟
    retry, err := d.SetnxEx(id, "1", 60 * 60)
    if err != nil {
        // 如果设置锁时发生错误，返回nil和false
        return nil, false
    }
	// 获取到锁，说明是第一次请求或者是时间过期了
    if retry {
        // 如果成功获取锁，返回nil和true，表示可以继续处理请求
        return nil, true
    }
    // 如果未获取到锁，则从缓存中获取id对应的数据，同时返回false，不继续后续请求
    resp, _ = d.Cache.Get(id)
    // 返回获取到的数据和false，表示未获取到锁
    return resp, false
}

// SaveResp 保存响应结果
// 该方法将给定的响应保存在缓存中，以便后续的相同请求可以直接返回该结果，而无需重新计算
func (d *defaultIdempotent) SaveResp(ctx context.Context, id string, resp interface{}, respErr error) error {
    // 保存结果
    d.Cache.Set(id, resp)
    return nil
}






































