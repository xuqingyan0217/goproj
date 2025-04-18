package job

import (
	"context"
	"errors"
	"time"
)

var ErrJobTimeout = errors.New("任务超时")

// RetryJetLagFunc 定义重试的时间策略
type RetryJetLagFunc func(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration

// RetryJetLagAlways 定义默认的方法
func RetryJetLagAlways(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration {
	return DefaultRetryJetLag
}

// IsRetryFunc 是否进行重试
type IsRetryFunc func(ctx context.Context, retryCount int, err error) bool

// IsRetryAlways 定义默认的方法
func IsRetryAlways(ctx context.Context, retryCount int, err error) bool {
	return true
}

// WithRetry 是一个用于处理具有重试逻辑的异步操作的函数。
// 它接受一个上下文、一个处理函数以及可变的重试选项作为参数。
// 其中，该处理函数看似一点逻辑没得，但事实上，使用它的时候，会往其中传入一个ctx，如果该ctx出了问题，那就会返回错误。
// 该函数会尝试执行处理函数，如果执行失败则根据重试策略进行重试。
func WithRetry(ctx context.Context, handler func(ctx context.Context) error, opts ...RetryOptions) error {
    // 依据输入参数，将其注入到option的初始化里
    opt := newOptions(opts...)

    // 从 ctx 上下文中获取截止时间，如果未设置，则返回false
    _, ok := ctx.Deadline()
	// 如下就是在未设置的情况下准备进行设置重试这一任务自己的超时时间（等价于if ok == false）
    if !ok {
		// 创建一个带有超时的Context，以控制后续操作不会无限期地执行下去。
		// 这里的目的是为了避免在某些情况下操作可能需要较长时间才能完成，从而导致程序挂起。
		// 通过使用context.WithTimeout，我们为操作设置了最大执行时间（opt.timeout），确保了程序的响应性和稳定性。
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.timeout)
		// 使用defer语句确保在函数返回前取消Context，释放相关资源，防止资源泄露。
		defer cancel()
    }
	// 到此处，说明ctx有截止时间
    var (
        herr error
        retryJetLag time.Duration
        ch = make(chan error, 1)
    )

    // 开始重试循环，尝试执行处理函数直到成功或达到最大重试次数。
    for i := 0; i < opt.retryNums; i++ {
        // 在goroutine中执行处理函数，以便主goroutine可以继续进行其他操作。
        go func() {
			// 如下是每次循环的时候将ctx放入到处理函数里，如果ctx时间截至了就报错。
            herr = handler(ctx)
			// 将错误传递给channel，便于后续select处理
            ch <- herr
        }()
        // 使用select语句来处理处理函数的返回值或上下文的完成。
        select {
        case herr = <-ch:
            // 如果处理函数成功执行（无错误），则退出函数，说明无需再重试了。
            if herr == nil {
                return nil
            }
            // 如果根据重试策略当前错误不应被重试，则返回该错误。
            if !opt.isRetryFunc(ctx, i, herr) {
                return herr
            }
            // 得到下一次重试前的等待时间，默认是1s，可随opt变化，并休眠当前goroutine。
            retryJetLag = opt.retryJetLag(ctx, i, retryJetLag)
            time.Sleep(retryJetLag)
        case <-ctx.Done():
            // 如果上下文在等待期间完成，则返回超时错误，到这说明重试这个任务本身也超时了。
            return ErrJobTimeout
        }
    }
    // 如果所有重试都失败，则返回最后一次尝试的错误。
    return herr
}

