package service

import (
	"context"
	"fmt"

	order "gomall/app/frontend/hertz_gen/frontend/order"
	"gomall/app/frontend/infra/rpc"
	rpcOrder "gomall/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type CancelOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCancelOrderService(Context context.Context, RequestContext *app.RequestContext) *CancelOrderService {
	return &CancelOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *CancelOrderService) Run(req *order.CancelOrderReq) (resp map[string]any, err error) {
	fmt.Println("-0-0-0-0-0", req)
	CancelResp, err := rpc.OrderClient.CancelOrder(h.Context, &rpcOrder.CancelOrderReq{
		OrderId: req.OrderId,
	})
	if err != nil {
		return nil, err
	}
	return utils.H{
		"success": CancelResp.Success,
	}, nil
}
