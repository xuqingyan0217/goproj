package middleware

import (
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

// LimitMiddleware 是一个中间件结构体，用于处理请求速率限制。
// 它结合了Redis配置和令牌桶限流算法。
type LimitMiddleware struct {
	redisCfg redis.RedisConf // redisCfg 存储Redis数据库配置信息。
	*limit.TokenLimiter      // TokenLimiter 嵌入的令牌限流器，用于执行实际的速率限制。
}

// NewLimitMiddleware 创建并返回一个新的LimitMiddleware实例。
// 参数cfg是Redis数据库的配置信息，用于初始化中间件。
func NewLimitMiddleware(cfg redis.RedisConf) *LimitMiddleware {
	return &LimitMiddleware{
		redisCfg:     cfg, // 初始化中间件时设置Redis配置信息。
	}
}

// TokenLimitHandler 返回一个REST中间件，用于应用令牌桶算法进行速率限制。
// 参数rate指定每秒生成的令牌数，参数burst指定令牌桶的容量。
// 该中间件使用Redis存储令牌桶状态，键名为"REDIS_TOKEN_LIMIT_KEY"。
func (l *LimitMiddleware) TokenLimitHandler(rate, burst int) rest.Middleware {
	l.TokenLimiter = limit.NewTokenLimiter(rate, burst, redis.MustNewRedis(l.redisCfg),"REDIS_TOKEN_LIMIT_KEY")

	// 返回的中间件函数，用于包裹HTTP处理函数，实现速率限制。
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 判断当前请求是否允许通过限流器。
			if l.TokenLimiter.AllowCtx(r.Context()) {
				// 如果允许，调用下一个处理函数。
				next(w, r)
				return
			}
			// 如果不允许，返回未授权状态码。
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

