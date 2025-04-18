package utils

import (
	"context"
)

// GetUserIdFormCtx 从上下文中获取用户ID
// 在JWT认证方式下，用户ID会通过GlobalAuth中间件设置到上下文中
func GetUserIdFormCtx(ctx context.Context) int32 {
	userId := ctx.Value(CtxUserID)
	if userId == nil {
		return 0
	}
	return userId.(int32)
}
