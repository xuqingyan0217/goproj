package utils

import (
	"github.com/golang-jwt/jwt/v4"
)

// ParseToken 将字符串token解析为*jwt.Token类型
// 参数:
//   - tokenString: 字符串类型的JWT token
// 返回:
//   - *jwt.Token: 解析后的token对象
//   - error: 解析过程中的错误
func ParseToken(tokenString string) (*jwt.Token, error) {
	// 使用与生成token相同的密钥进行解析
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		
		// 返回用于验证的密钥
		return []byte(JWTSecretKey), nil
	})
	
	return token, err
}

// GetUserIDFromToken 从token字符串中提取用户ID
// 参数:
//   - tokenString: 字符串类型的JWT token
// 返回:
//   - int32: 用户ID
//   - error: 解析过程中的错误
func GetUserIDFromToken(tokenString string) (int32, error) {
	// 解析token
	token, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	
	// 验证token是否有效
	if !token.Valid {
		return 0, jwt.NewValidationError("invalid token", jwt.ValidationErrorUnverifiable)
	}
	
	// 提取claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.NewValidationError("invalid claims", jwt.ValidationErrorClaimsInvalid)
	}
	
	// 提取用户ID
	userIDFloat, ok := claims[TokenUserId].(float64)
	if !ok {
		return 0, jwt.NewValidationError("user id not found in token", jwt.ValidationErrorClaimsInvalid)
	}
	
	return int32(userIDFloat), nil
}