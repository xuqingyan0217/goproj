package service

import (
	"context"
	"gomall/app/order/biz/dal/model"
	"gomall/app/order/biz/dal/mysql"
	order "gomall/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

type CancelPaymentService struct {
	ctx context.Context
} // NewCancelPaymentService new CancelPaymentService
func NewCancelPaymentService(ctx context.Context) *CancelPaymentService {
	return &CancelPaymentService{ctx: ctx}
}

// Run create note info
func (s *CancelPaymentService) Run(req *order.CancelPaymentReq) (resp *order.CancelPaymentResp, err error) {
	// 更新订单的支付状态
	err = mysql.DB.Model(&model.Order{}).Where("order_id = ?", req.OrderId).Update("user_currency", req.Status).Error
	if err != nil {
		return nil, kerrors.NewBizStatusError(500001, err.Error())
	}

	return &order.CancelPaymentResp{
		Message: "cancel payment success",
		Success: true,
	}, nil
}
