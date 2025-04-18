package AIEino

import (
	"context"
	AIEino "gomall/rpc_gen/kitex_gen/AIEino"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func AIWithOrders(ctx context.Context, req *AIEino.AIWithOrdersReq, callOptions ...callopt.Option) (resp *AIEino.AIWithOrdersResp, err error) {
	resp, err = defaultClient.AIWithOrders(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AIWithOrders call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AIWithPreCheckout(ctx context.Context, req *AIEino.AIWithPreCheckoutReq, callOptions ...callopt.Option) (resp *AIEino.AIWithPreCheckoutResp, err error) {
	resp, err = defaultClient.AIWithPreCheckout(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AIWithPreCheckout call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
