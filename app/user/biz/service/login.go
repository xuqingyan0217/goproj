package service

import (
	"context"
	"errors"
	"gomall/app/user/biz/dal/mysql"
	"gomall/app/user/biz/model"
	user "gomall/rpc_gen/kitex_gen/user"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginService struct {
	ctx context.Context
} // NewLoginService new LoginService
func NewLoginService(ctx context.Context) *LoginService {
	return &LoginService{ctx: ctx}
}

// Run create note info
func (s *LoginService) Run(req *user.LoginReq) (resp *user.LoginResp, err error) {
	// Finish your business logic.
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("param error")
	}
	row, err := model.GetByEmail(mysql.DB, req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(row.PasswordHashed), []byte(req.Password)); err != nil {
		return nil, errors.New("password not right")
	}
	resp = &user.LoginResp{
		UserId: int32(row.ID),
	}
	return resp, nil
}
