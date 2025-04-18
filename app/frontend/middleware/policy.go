package middleware

// RolePermissions 定义角色权限映射
var RoleProductPermissions = map[string][]string{
	"admin": {
		"/product/api/*",
	},
	"user": {
		
	},
}
// 后续针对更多路由组的角色认证，都可以新增一个map

// WhiteList 白名单路径，允许所有用户访问
var WhiteList = []string{
	"/api",
	"/sign-in",
	"/sign-up",
	"/product",
}

// BlackList 黑名单路径，禁止所有用户访问
var BlackList = []string{
	"/internal/*",
	"/admin/system/*",
	"/api/debug/*",
}

// UserBlackList 用户黑名单，存储被禁止访问的用户ID
var UserBlackList = map[int32]bool{
	3:  true, // 被封禁的用户ID
	10: true, // 恶意用户
}
