package service

import (
	"context"
	"errors"
	"regexp"
	"golang.org/x/crypto/bcrypt"
	"gomall/app/user/biz/dal/mysql"
	"gomall/app/user/biz/model"
	user "gomall/rpc_gen/kitex_gen/user"
)

type RegisterService struct {
	ctx context.Context
} // NewRegisterService new RegisterService
func NewRegisterService(ctx context.Context) *RegisterService {
	return &RegisterService{ctx: ctx}
}

// Run create note info
func (s *RegisterService) Run(req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// 验证参数
	if req.Email == "" || req.Password == "" || req.PasswordConfirm == "" {
		return nil, errors.New("param error")
	}

	// 验证邮箱格式
	emailLen := len(req.Email)
	if emailLen < 5 || emailLen > 100 {
		return nil, errors.New("email length must be between 5 and 100 characters")
	}
	emailRegex := `^[\w-]+(\.[\w-]+)*@[\w-]+(\.[\w-]+)+$`
	if !regexp.MustCompile(emailRegex).MatchString(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// 检查邮箱是否已被注册
	existingUser, err := model.GetByEmail(mysql.DB, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// 验证密码格式
	passwordLen := len(req.Password)
	if passwordLen < 8 || passwordLen > 20 {
		return nil, errors.New("password length must be between 8 and 20 characters")
	}
	passwordRegex := `^[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*\d+[a-zA-Z\d]*$|^[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*\d+[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*$|^[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*\d+[a-zA-Z\d]*$|^[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*\d+[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*$|^[a-zA-Z\d]*\d+[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*$|^[a-zA-Z\d]*\d+[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*$`
	if !regexp.MustCompile(passwordRegex).MatchString(req.Password) {
		return nil, errors.New("password must contain uppercase, lowercase letters and numbers")
	}

	// 验证密码确认
	if req.Password != req.PasswordConfirm {
		return nil, errors.New("password not match")
	}
	// 加密
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser := &model.User{
		Email:          req.Email,
		PasswordHashed: string(passwordHashed),
	}
	err = model.Create(mysql.DB, newUser)
	if err != nil {
		return nil, err
	}
	return &user.RegisterResp{
		UserId: int32(newUser.ID),
	}, nil
}
