package rpcserver

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func SyncXLimitInterceptor(maxCount int) grpc.UnaryServerInterceptor {
	// 声明一个人限流对象
	l := syncx.NewLimit(maxCount)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 判断是否限流，if里是没有限流，直接走handler
		if l.TryBorrow() {
			defer func() {
				if err := l.Return(); err != nil {
					logx.Error(err)
				}
			}()
			return handler(ctx, req)
		} else {
			// 若限流，则抛出错误
			logx.Errorf("concurrent connections over %d, rejected with code %d", maxCount, http.StatusServiceUnavailable)
			return nil, status.Error(codes.Unavailable, "concurrent connections over limit")
		}

	}
}
