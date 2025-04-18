package middleware

import (
	"easy-chat/pkg/interceptor"
	"net/http"
)

type IdempotenceMiddleware struct {

}

func NewIdempotenceMiddleware() *IdempotenceMiddleware {
	return &IdempotenceMiddleware{}
}

// 中间件的处理很简单，就是把请求先到这里，给请求头加个value，然后继续正常的请求

// Handler 是一个中间件处理函数，用于在请求处理中加入幂等性控制。
// 它接收一个 http.HandlerFunc 类型的 next 参数，表示当前中间件之后的处理函数。
// 返回一个新的 http.HandlerFunc，它会在调用 next 处理函数之前执行额外的幂等性控制逻辑。
func (i *IdempotenceMiddleware) Handler(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 在请求的上下文中添加幂等性控制相关的值。
        // 这里使用了 interceptor.ContextWithVal 来设置上下文，这可能包括幂等性校验所需的标识或状态。
        r = r.WithContext(interceptor.ContextWithVal(r.Context()))

        // 调用下一个处理函数，保持请求处理的链式调用。
        // 这确保了在幂等性控制逻辑执行后，请求能继续被后续的处理函数处理。
        next(w, r)
    }
}


