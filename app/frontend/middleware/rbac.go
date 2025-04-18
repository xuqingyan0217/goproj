package middleware

import (
	"context"
	frontentUtils "gomall/app/frontend/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// RBAC 基于角色的访问控制中间件
func RBAC(MyPolicy map[string][]string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从JWT claims中获取用户角色
		claims, err := frontentUtils.JWTMiddleware.GetClaimsFromJWT(ctx, c)
		if err != nil {
			c.JSON(401, map[string]interface{}{
				"code":    401,
				"message": "未登录或登录已过期",
				"data": map[string]interface{}{
					"redirect": "/sign-in?next=" + c.FullPath(),
				},
			})
			c.Abort()
			return
		}
		if claims == nil {
			c.JSON(403, map[string]interface{}{
				"code":    403,
				"message": "无法获取用户角色信息",
			})
			c.Abort()
			return
		}

		// 获取用户ID
		userID := int32(claims[frontentUtils.TokenUserId].(float64))
		// 检查用户是否在黑名单中
		if UserBlackList[userID] {
			c.JSON(403, map[string]interface{}{
				"code":    403,
				"message": "您的账号已被禁止访问",
			})
			c.Abort()
			return
		}

		// 获取用户角色
		role, ok := claims[frontentUtils.TokenUserRole].(string)
		if !ok {
			c.JSON(403, map[string]interface{}{
				"code":    403,
				"message": "用户角色信息无效",
			})
			c.Abort()
			return
		}

		// 获取当前请求路径
		path := string(c.Request.URI().Path())

		// 检查用户角色是否有权限访问该路径
		if !hasPermission(role, path, MyPolicy) {
			c.JSON(403, map[string]interface{}{
				"code":    403,
				"message": "没有访问权限",
			})
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// hasPermission 检查指定角色是否有权限访问指定路径
func hasPermission(role string, path string, policy map[string][]string) bool {
	// 首先检查是否在黑名单中
	for _, pattern := range BlackList {
		if matchPath(pattern, path) {
			return false
		}
	}

	// 然后检查是否在白名单中
	for _, pattern := range WhiteList {
		if matchPath(pattern, path) {
			return true
		}
	}

	// 最后检查角色权限
	permissions, exists := policy[role]
	if !exists {
		return false
	}

	// 检查路径是否在权限列表中
	for _, pattern := range permissions {
		if matchPath(pattern, path) {
			return true
		}
	}

	return false
}

// matchPath 检查请求路径是否匹配权限模式
func matchPath(pattern string, path string) bool {
	// 如果模式是通配符*，允许访问所有路径
	if pattern == "*" {
		return true
	}

	// 精确匹配路径
	if pattern == path {
		return true
	}

	// 处理通配符情况
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(path) >= len(prefix) && path[:len(prefix)] == prefix
	}

	return false
}
