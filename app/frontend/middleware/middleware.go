package middleware

import "github.com/cloudwego/hertz/pkg/app/server"

func Register(h *server.Hertz) {
	// 注册全局认证中间件，用于从token中提取用户ID
	h.Use(GlobalAuth())
}
