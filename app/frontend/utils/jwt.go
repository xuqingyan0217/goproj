package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
)

// JWTMiddleware 全局JWT中间件实例
var JWTMiddleware *jwt.HertzJWTMiddleware

// JWT密钥
const JWTSecretKey = "gomall_secret_key"

// TokenUserRole 用户角色的key
const TokenUserRole = "role"

// InitJWT 初始化JWT中间件
func InitJWT() {
	middleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       "gomall",
		Key:         []byte(JWTSecretKey),
		Timeout:     5 * time.Minute, // token有效期
		MaxRefresh:  5 * time.Minute, // token刷新有效期
		IdentityKey: TokenUserId,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if claims, ok := data.(map[string]interface{}); ok {
				userID := claims["user_id"].(int32)
				role := claims["role"].(string)
				return jwt.MapClaims{
					TokenUserId: float64(userID),
					TokenUserRole: role,
				}
			} else if v, ok := data.(int32); ok {
				// 向后兼容：如果只传入userID，默认角色为user
				return jwt.MapClaims{
					TokenUserId: float64(v),
					TokenUserRole: "user",
				}
			}
			// 如果data类型不正确，记录日志并返回空claims
			fmt.Println("Invalid data type in PayloadFunc:", data)
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			userID, _ := claims[TokenUserId].(float64)
			return int32(userID)
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(401, map[string]interface{}{
				"code":    401,
				"message": "Unauthorized",
				"data": map[string]interface{}{
					"redirect": "/sign-in?next=" + c.FullPath(),
				},
			})
			c.Abort()
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		panic("JWT初始化失败: " + err.Error())
	}

	// 初始化中间件
	err = middleware.MiddlewareInit()
	if err != nil {
		panic("JWT中间件初始化失败: " + err.Error())
	}

	JWTMiddleware = middleware
}

// GenerateToken 生成JWT token并验证
func GenerateToken(userID int32, role string) (string, time.Time, error) {
	claims := map[string]interface{}{
		"user_id": userID,
		"role":    role,
	}
	token, expire, err := JWTMiddleware.TokenGenerator(claims)
	if err != nil {
		return "", time.Time{}, err
	}

	// 解析token以验证userId
	userId, err := GetUserIDFromToken(token)
	if err!= nil {
		return "", time.Time{}, err
	}
	// 验证userId是否匹配
	if userId!= userID {
		return "", time.Time{}, fmt.Errorf("userId不匹配")
	}
	
	return token, expire, nil
}
