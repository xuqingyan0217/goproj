package service

import (
	"context"
	"gomall/app/order/biz/dal/model"
	"gomall/app/order/biz/dal/mysql"
	order "gomall/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

type ChangeOrderStatusService struct {
	ctx context.Context
} // NewChangeOrderStatusService new ChangeOrderStatusService
func NewChangeOrderStatusService(ctx context.Context) *ChangeOrderStatusService {
	return &ChangeOrderStatusService{ctx: ctx}
}

// Run create note info
func (s *ChangeOrderStatusService) Run(req *order.ChangeOrderStatusReq) (resp *order.ChangeOrderStatusResp, err error) {
	// Finish your business logic.
	// 更新订单的支付状态
	err = mysql.DB.Model(&model.Order{}).Where("order_id = ?", req.OrderId).Update("user_currency", req.Status).Error
	if err != nil {
		return nil, kerrors.NewBizStatusError(500001, err.Error())
	}
	return &order.ChangeOrderStatusResp{
		Success: true,
	}, nil
}
