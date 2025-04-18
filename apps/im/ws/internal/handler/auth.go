/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package handler

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"

	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/pkg/ctxdata"
	"net/http"
)

type JwtAuth struct {
	svc    *svc.ServiceContext
	parser *token.TokenParser
	logx.Logger
}

func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}

// Auth 执行 JWT 认证。
// 该方法尝试从请求中解析 JWT 令牌，验证其有效性，并获取用户声明。如果成功，它会将声明添加到请求上下文中，
// 表示认证通过，并返回 true；否则返回 false。
// 参数：
//   w - 用于写响应的响应写入器。
//   r - 指向传入请求的指针。
// 返回值：
//   返回一个布尔值，表示认证是否通过。
func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {

	if token := r.Header.Get("sec-websocket-protocol"); token != "" {
		r.Header.Set("Authorization", token)
	}

    // 尝试解析令牌，如果令牌无效或缺失可能会引发错误。
    tok, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
    if err != nil {
        // 记录令牌解析错误，用于调试，并返回 false 表示认证失败。
        j.Errorf("parse token err %v ", err)
        return false
    }

    // 检查令牌是否有效，如果无效则返回 false 表示认证失败。
    if !tok.Valid {
        return false
    }

    // 尝试将令牌声明转换为 jwt.MapClaims 类型，如果转换失败则返回 false。
    claims, ok := tok.Claims.(jwt.MapClaims)
    if !ok {
        return false
    }

    // 将声明添加到请求上下文中，表示认证成功并传递控制权给后续处理。
    *r = *r.WithContext(context.WithValue(r.Context(), ctxdata.Identify, claims[ctxdata.Identify]))

    // 返回 true 表示认证成功。
    return true
}


func (j *JwtAuth) UserId(r *http.Request) string {
	return ctxdata.GetUId(r.Context())
}
