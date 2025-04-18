package order

import (
	"context"
	order "gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func PlaceOrder(ctx context.Context, req *order.PlaceOrderReq, callOptions ...callopt.Option) (resp *order.PlaceOrderResp, err error) {
	resp, err = defaultClient.PlaceOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "PlaceOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListOrder(ctx context.Context, req *order.ListOrderReq, callOptions ...callopt.Option) (resp *order.ListOrderResp, err error) {
	resp, err = defaultClient.ListOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func CancelOrder(ctx context.Context, req *order.CancelOrderReq, callOptions ...callopt.Option) (resp *order.CancelOrderResp, err error) {
	resp, err = defaultClient.CancelOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CancelOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func CancelPayment(ctx context.Context, req *order.CancelPaymentReq, callOptions ...callopt.Option) (resp *order.CancelPaymentResp, err error) {
	resp, err = defaultClient.CancelPayment(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CancelPayment call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ChangeOrderStatus(ctx context.Context, req *order.ChangeOrderStatusReq, callOptions ...callopt.Option) (resp *order.ChangeOrderStatusResp, err error) {
	resp, err = defaultClient.ChangeOrderStatus(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ChangeOrderStatus call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
