package service

import (
	"context"
	"fmt"
	"gomall/app/order/biz/dal/model"
	"gomall/app/order/biz/dal/mysql"
	order "gomall/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

type CancelOrderService struct {
	ctx context.Context
} // NewCancelOrderService new CancelOrderService
func NewCancelOrderService(ctx context.Context) *CancelOrderService {
	return &CancelOrderService{ctx: ctx}
}

// Run create note info
func (s *CancelOrderService) Run(req *order.CancelOrderReq) (resp *order.CancelOrderResp, err error) {
	// 开启事务
	fmt.Println("-0-0-0-0", req)
	tx := mysql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 更新订单的状态
	err = mysql.DB.Model(&model.Order{}).Where("order_id = ?", req.OrderId).Update("user_currency", "CANCEL").Error
	if err != nil {
		return nil, kerrors.NewBizStatusError(500001, err.Error())
	}
	// 先删除订单关联的订单项
	if err = tx.Where("order_id_refer = ?", req.OrderId).Delete(&model.OrderItem{}).Error; err != nil {
		tx.Rollback()
		return nil, kerrors.NewBizStatusError(500001, err.Error())
	}

	// 再删除订单本身
	if err = tx.Where("order_id = ?", req.OrderId).Delete(&model.Order{}).Error; err != nil {
		tx.Rollback()
		return nil, kerrors.NewBizStatusError(500001, err.Error())
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, kerrors.NewBizStatusError(500001, err.Error())
	}

	return &order.CancelOrderResp{
		Message: "cancel order success",
		Success: true,
	}, nil
}
