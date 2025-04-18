package logic

import (
	"context"
	"easy-chat/pkg/ctxdata"
	"regexp"
	"time"

	"github.com/pkg/errors"

	"easy-chat/pkg/xerr"

	"easy-chat/apps/user/models"
	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"
	"easy-chat/pkg/encrypt"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号没有注册")
	ErrUserPwdError     = xerr.New(xerr.SERVER_COMMON_ERROR, "密码不正确")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// todo: add your logic here and delete this line

	// 判断用户的输入手机号格式
	pattern := `^(13\d|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18\d|19[0-35-9])\d{8}$`
	// 正则校验
	regexpPattern := regexp.MustCompile(pattern)
	if !regexpPattern.MatchString(in.Phone) {
		return nil, errors.New("手机号格式不正确")
	}

	// 判断之前设置的token是否过期
	// 过期则触发获取验证码，然后重新登录
	// 没过期就下一步

	// 1. 验证用户是否注册，根据手机号码验证
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v , req %v", err, in.Phone)
	}

	// 密码验证
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, errors.WithStack(ErrUserPwdError)
	}

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire,
		userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err %v", err)
	}

	//return nil, errors.New("做测试")
	return &user.LoginResp{
		// Id: userEntity.Id,
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
		User: &user.UserEntity{
			Id:       userEntity.Id,
			Phone:    in.Phone,
			Nickname: userEntity.Nickname,
			Avatar:   userEntity.Avatar,
			Sex:      int32(userEntity.Sex.Int64),
		},
	}, nil
}
