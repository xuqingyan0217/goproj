package job

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithRetry(t *testing.T) {
	var (
		ErrTest = errors.New("测试异常")
		handler = func(ctx context.Context) error {
			t.Log("执行任务")
			return ErrTest
		}
	)
	type args struct {
		ctx     context.Context
		handler func(context.Context) error
		opts    []RetryOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			// 会任务超时，重试次数执行不完
			"1", args{
			ctx: context.Background(),
			handler: handler,
			opts:    []RetryOptions{},
		}, ErrJobTimeout,
		},
		{
			// 能成功执行完所有超时，因为时间长了，执行一次后的时间短了
			"2", args{
			ctx: context.Background(),
			handler: handler,
			opts:    []RetryOptions{
				WithRetryTime(3 * time.Second),
				WithRetryJetLag(func(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration{
					return 500 * time.Millisecond
				}),
			},
		}, ErrTest,
		},
		{
			// 静止重试，也就是执行一次后就停止，即便该次有错误需要进行重试
			"3", args{
			ctx: context.Background(),
			handler: handler,
			opts:    []RetryOptions{
				WithIsRetryFunc(func(ctx context.Context, retryCount int, err error) bool {
					return false
				}),
			},
		}, ErrTest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WithRetry(tt.args.ctx, tt.args.handler, tt.args.opts...); err != tt.wantErr {
				t.Errorf("WithRetry() err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}