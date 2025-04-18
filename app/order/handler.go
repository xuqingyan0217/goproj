package main

import (
	"context"
	order "gomall/rpc_gen/kitex_gen/order"
	"gomall/app/order/biz/service"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// PlaceOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	resp, err = service.NewPlaceOrderService(ctx).Run(req)

	return resp, err
}

// ListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	resp, err = service.NewListOrderService(ctx).Run(req)

	return resp, err
}

// CancelOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, req *order.CancelOrderReq) (resp *order.CancelOrderResp, err error) {
	resp, err = service.NewCancelOrderService(ctx).Run(req)

	return resp, err
}

// CancelPayment implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CancelPayment(ctx context.Context, req *order.CancelPaymentReq) (resp *order.CancelPaymentResp, err error) {
	resp, err = service.NewCancelPaymentService(ctx).Run(req)

	return resp, err
}

// ChangeOrderStatus implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ChangeOrderStatus(ctx context.Context, req *order.ChangeOrderStatusReq) (resp *order.ChangeOrderStatusResp, err error) {
	resp, err = service.NewChangeOrderStatusService(ctx).Run(req)

	return resp, err
}
