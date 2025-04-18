package service

import (
	"context"
	"fmt"

	common "gomall/app/frontend/hertz_gen/frontend/common"
	"gomall/app/frontend/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

type RefreshTokenService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRefreshTokenService(Context context.Context, RequestContext *app.RequestContext) *RefreshTokenService {
	return &RefreshTokenService{RequestContext: RequestContext, Context: Context}
}

func (h *RefreshTokenService) Run(req *common.Empty) (resp *common.Empty, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()

	// 从请求头中获取当前令牌
	tokenString := h.RequestContext.GetHeader("Authorization")
	if len(tokenString) == 0 {
		return nil, fmt.Errorf("未提供令牌")
	}
	// 如果令牌带有Bearer前缀，则去除
	if len(tokenString) > 7 && string(tokenString[:7]) == "Bearer " {
		tokenString = tokenString[7:]
	}

	// 验证令牌并提取用户ID
	userID, err := utils.GetUserIDFromToken(string(tokenString))
	if err != nil {
		return nil, fmt.Errorf("令牌验证失败: %w", err)
	}
	role := "user"
	// 生成JWT token
	if userID == 1 {
		role = "admin"
	}
	// 生成新的令牌
	newToken, _, err := utils.GenerateToken(userID, role)
	if err != nil {
		return nil, fmt.Errorf("生成新令牌失败: %w", err)
	}

	// 设置新令牌到响应头
	h.RequestContext.Header("Authorization", "Bearer "+newToken)

	return &common.Empty{}, nil
}
