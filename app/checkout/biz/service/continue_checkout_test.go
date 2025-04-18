package service

import (
	"context"
	"testing"
	checkout "gomall/rpc_gen/kitex_gen/checkout"
)

func TestContinueCheckout_Run(t *testing.T) {
	ctx := context.Background()
	s := NewContinueCheckoutService(ctx)
	// init req and assert value

	req := &checkout.ContinueCheckoutReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
