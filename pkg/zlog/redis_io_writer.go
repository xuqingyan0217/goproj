package zlog

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"io"
)

type redisIoWriter struct {
	redisKey string
	redis *redis.Redis
}

func NewRedisIoWriter(redisKey string, redisCfg redis.RedisConf) io.Writer {
	return &redisIoWriter{
		redisKey: redisKey,
		redis:    redis.MustNewRedis(redisCfg),
	}
}

func (r *redisIoWriter) Write(p []byte) (n int, err error) {
	// 将日志写入redis
	go r.redis.Rpush(r.redisKey, string(p))

	fmt.Println("已经将日志写入redis" ,r.redisKey,string(p))

	return n, nil
}

