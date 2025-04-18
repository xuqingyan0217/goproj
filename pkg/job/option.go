package job

import "time"

type (
	RetryOptions func(opt *retryOptions)

	retryOptions struct {
		timeout time.Duration
		retryNums int
		isRetryFunc IsRetryFunc
		retryJetLag RetryJetLagFunc
	}
)

func newOptions(opts ...RetryOptions) *retryOptions {
	opt := &retryOptions{
		timeout:     DefaultRetryTimeout,
		retryNums:   DefaultRetryNums,
		isRetryFunc: IsRetryAlways,
		retryJetLag: RetryJetLagAlways,
	}
	for _, options := range opts {
		options(opt)
	}
	return opt
}

func WithRetryTime(timeout time.Duration) RetryOptions {
	return func(opt *retryOptions) {
		if timeout > 0 {
			opt.timeout = timeout
		}
	}
}

func WithRetryNums(retryNums int) RetryOptions {
	return func(opt *retryOptions) {
		if retryNums > 1 {
			opt.retryNums = retryNums
		}
	}
}

func WithIsRetryFunc(retryFunc IsRetryFunc) RetryOptions {
	return func(opt *retryOptions) {
		if retryFunc != nil {
			opt.isRetryFunc = retryFunc
		}
	}
}

func WithRetryJetLag(retryJetLagFunc RetryJetLagFunc) RetryOptions {
	return func(opt *retryOptions) {
		if retryJetLagFunc != nil {
			opt.retryJetLag = retryJetLagFunc
		}
	}
}

