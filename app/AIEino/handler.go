package main

import (
	"context"
	"gomall/app/AIEino/biz/service"
	AIEino "gomall/rpc_gen/kitex_gen/AIEino"
)

// AIEinoServiceImpl implements the last service interface defined in the IDL.
type AIEinoServiceImpl struct{}

// AIWithOrders implements the AIEinoServiceImpl interface.
func (s *AIEinoServiceImpl) AIWithOrders(ctx context.Context, req *AIEino.AIWithOrdersReq) (resp *AIEino.AIWithOrdersResp, err error) {
	resp, err = service.NewAIWithOrdersService(ctx).Run(req)

	return resp, err
}

// AIWithPreCheckout implements the AIEinoServiceImpl interface.
func (s *AIEinoServiceImpl) AIWithPreCheckout(ctx context.Context, req *AIEino.AIWithPreCheckoutReq) (resp *AIEino.AIWithPreCheckoutResp, err error) {
	resp, err = service.NewAIWithPreCheckoutService(ctx).Run(req)

	return resp, err
}
