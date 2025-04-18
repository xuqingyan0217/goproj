package service

import (
	"context"
	"gomall/app/frontend/infra/rpc"
	"gomall/app/frontend/utils"
	"gomall/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
	auth "gomall/app/frontend/hertz_gen/frontend/auth"
)

type LoginService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLoginService(Context context.Context, RequestContext *app.RequestContext) *LoginService {
	return &LoginService{RequestContext: RequestContext, Context: Context}
}

func (h *LoginService) Run(req *auth.LoginReq) (redirect string, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	resp, err := rpc.UserClient.Login(h.Context, &user.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return "", err
	}
	role := "user"
	// 生成JWT token
	if resp.UserId == 1 {
		role = "admin"
	}
	token, _, err := utils.GenerateToken(resp.UserId, role)
	if err != nil {
		return "", err
	}
	
	// 设置token到响应头
	h.RequestContext.Header("Authorization", "Bearer "+token)
	
	redirect = "/"
	if req.Next != "" {
		redirect = req.Next
	}
	return
}
