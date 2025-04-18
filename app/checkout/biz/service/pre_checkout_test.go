package service

import (
	"context"
	"testing"
	checkout "gomall/rpc_gen/kitex_gen/checkout"
)

func TestPreCheckout_Run(t *testing.T) {
	ctx := context.Background()
	s := NewPreCheckoutService(ctx)
	// init req and assert value

	req := &checkout.PreCheckoutReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
