package checkout

import (
	"context"
	checkout "gomall/rpc_gen/kitex_gen/checkout"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func Checkout(ctx context.Context, req *checkout.CheckoutReq, callOptions ...callopt.Option) (resp *checkout.CheckoutResp, err error) {
	resp, err = defaultClient.Checkout(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Checkout call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func PreCheckout(ctx context.Context, req *checkout.PreCheckoutReq, callOptions ...callopt.Option) (resp *checkout.PreCheckoutResp, err error) {
	resp, err = defaultClient.PreCheckout(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "PreCheckout call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ContinueCheckout(ctx context.Context, req *checkout.ContinueCheckoutReq, callOptions ...callopt.Option) (resp *checkout.ContinueCheckoutResp, err error) {
	resp, err = defaultClient.ContinueCheckout(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ContinueCheckout call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
