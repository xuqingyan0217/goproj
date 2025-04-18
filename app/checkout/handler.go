package main

import (
	"context"
	checkout "gomall/rpc_gen/kitex_gen/checkout"
	"gomall/app/checkout/biz/service"
)

// CheckoutServiceImpl implements the last service interface defined in the IDL.
type CheckoutServiceImpl struct{}

// Checkout implements the CheckoutServiceImpl interface.
func (s *CheckoutServiceImpl) Checkout(ctx context.Context, req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	resp, err = service.NewCheckoutService(ctx).Run(req)

	return resp, err
}

// PreCheckout implements the CheckoutServiceImpl interface.
func (s *CheckoutServiceImpl) PreCheckout(ctx context.Context, req *checkout.PreCheckoutReq) (resp *checkout.PreCheckoutResp, err error) {
	resp, err = service.NewPreCheckoutService(ctx).Run(req)

	return resp, err
}

// ContinueCheckout implements the CheckoutServiceImpl interface.
func (s *CheckoutServiceImpl) ContinueCheckout(ctx context.Context, req *checkout.ContinueCheckoutReq) (resp *checkout.ContinueCheckoutResp, err error) {
	resp, err = service.NewContinueCheckoutService(ctx).Run(req)

	return resp, err
}
