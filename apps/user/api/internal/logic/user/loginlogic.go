package user

import (
	"context"
	"easy-chat/apps/user/rpc/user"
	"easy-chat/pkg/constants"
	"github.com/jinzhu/copier"

	"easy-chat/apps/user/api/internal/svc"
	"easy-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	// fmt.Println(l.svcCtx.Config.Database)

	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	var res types.LoginResp
	copier.Copy(&res, loginResp)

	// 处理调用rpc登录时，添加用户信息到cache的业务
	// 当用户登入时，在hash里面，将value设置为1，表示在线
	l.svcCtx.Redis.HsetCtx(l.ctx, constants.REDIS_OLINE_USER, loginResp.Id, "1")

	return &res, nil
}
