package service

import (
	"context"
	"testing"
	order "gomall/rpc_gen/kitex_gen/order"
)

func TestChangeOrderStatus_Run(t *testing.T) {
	ctx := context.Background()
	s := NewChangeOrderStatusService(ctx)
	// init req and assert value

	req := &order.ChangeOrderStatusReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
