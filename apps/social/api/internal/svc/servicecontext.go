package svc

import (
	"easy-chat/apps/im/rpc/imclient"
	"easy-chat/apps/social/api/internal/config"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/user/rpc/userclient"
	"easy-chat/pkg/interceptor"
	"easy-chat/pkg/interceptor/rpcclient"
	"easy-chat/pkg/middleware"
	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"time"
)

var retryPolicy = `{
	"methodConfig": [{
		"name": [{
			"service": "social.social"
		}],
		"waitForReady": true,
		"retryPolicy": {
			"maxAttempts": 5,
			"initialBackoff": "0.001s",
			"maxBackoff": "0.002s",
			"backoffMultiplier": 1.0,
			"retryableStatusCodes": ["UNKNOWN", "DEADLINE_EXCEEDED"]
		}
	}]
}`

type ServiceContext struct {
	Config config.Config
	// 添加中间件
	IdempotenceMiddleware rest.Middleware
	LimitMiddleware rest.Middleware
	*redis.Redis
	socialclient.Social
	userclient.User
	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc,
			// 重试机制
			zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)),
			// 客户端拦截器
			zrpc.WithUnaryClientInterceptor(interceptor.DefaultIdempotentClient))),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc,
			zrpc.WithUnaryClientInterceptor(rpcclient.NewSheddingClient("user-rpc",
				load.WithBuckets(10),
				load.WithWindow(time.Millisecond * 10000),
				load.WithCpuThreshold(100),
				)),
		)),
		Im:	imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		Redis: redis.MustNewRedis(c.Redisx),
		// 添加中间件
		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware().Handler,
		LimitMiddleware: middleware.NewLimitMiddleware(c.Redisx).TokenLimitHandler(1, 30),
	}
}
