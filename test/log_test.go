package main

import (
	"context"
	"easy-chat/pkg/zlog"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/trace"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"testing"
)

func TestFile_log(t *testing.T) {
	logx.SetUp(logx.LogConf{
		ServiceName: "file.log",
		Mode:        "file",
		Encoding:    "json",
		Path:        "./z_log",
	})

	logx.Info("test file log")
	logx.Info("这是测试go-zero的日志")
	logx.Info("日志将由filebeat推送给elk")
	for {
	}
}

func TestRedis_log_ioWriter(t *testing.T) {
	io := zlog.NewRedisIoWriter("redis.io.write", redis.RedisConf{
		Host:        "10.0.0.17:16379",
		Type:        "node",
		Pass:        "easy-chat",
		Tls:         false,
		NonBlock:    false,
		PingTimeout: 0,
	})

	logx.SetWriter(logx.NewWriter(io))

	logx.Info("test file log by redis")
	logx.Info("这是测试go-zero-redis的日志")
	logx.Info("日志将由go-zero推送给redis")
	logx.Infow("带有key-value标识的日志", logx.LogField{
		Key:   "rid",
		Value: "11111",
	})
	for {
	}
}

func TestCtx_Log(t *testing.T) {
	io := zlog.NewRedisIoWriter("redis.io.write", redis.RedisConf{
		Host:        "10.0.0.17:16379",
		Type:        "node",
		Pass:        "easy-chat",
		Tls:         false,
		NonBlock:    false,
		PingTimeout: 0,
	})

	logx.SetWriter(logx.NewWriter(io))

	ctx, _ := sdkTrace.NewTracerProvider().Tracer(trace.TraceName).Start(context.Background(), "a")

	log := logx.WithContext(ctx)

	log.Info("test file log by ctx")
	log.Info("这是测试go-zero-ctx的日志")
	log.Info("日志将由go-zero通过ctx获取")
	fmt.Println("------------------------------------")
	logx.Info("test file log by redis")
	logx.Info("这是测试go-zero-redis的日志")
	logx.Info("日志将由go-zero推送给redis")
	for {
	}
}
