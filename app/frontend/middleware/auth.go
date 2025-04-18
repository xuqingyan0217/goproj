package middleware

import (
	"context"
	"fmt"
	frontentUtils "gomall/app/frontend/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

func GlobalAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查Authorization头的格式
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 0 {
			// 确保前缀正确 - 应该是"Bearer "而不是其他格式
			if len(authHeader) > 7 && string(authHeader[0:7]) == "Bearer " {
				fmt.Println("Bearer prefix found correctly")

				// 直接从Authorization头中提取token字符串
				tokenString := string(authHeader[7:])
				fmt.Println("Extracted token:", tokenString)

				// 使用utils包中的ParseToken函数解析token
				userId, err := frontentUtils.GetUserIDFromToken(tokenString)
				if err != nil {
					fmt.Println("Failed to parse token:", err)
				} else {
					ctx = context.WithValue(ctx, frontentUtils.CtxUserID, int32(userId))
				}
			}
		}
		c.Next(ctx)
	}
}

func Auth() app.HandlerFunc {
	return frontentUtils.JWTMiddleware.MiddlewareFunc()
}
