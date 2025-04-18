package AIEino

import (
	"context"
	AIEino "gomall/rpc_gen/kitex_gen/AIEino"

	"gomall/rpc_gen/kitex_gen/AIEino/aieinoservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() aieinoservice.Client
	Service() string
	AIWithOrders(ctx context.Context, Req *AIEino.AIWithOrdersReq, callOptions ...callopt.Option) (r *AIEino.AIWithOrdersResp, err error)
	AIWithPreCheckout(ctx context.Context, Req *AIEino.AIWithPreCheckoutReq, callOptions ...callopt.Option) (r *AIEino.AIWithPreCheckoutResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := aieinoservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient aieinoservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() aieinoservice.Client {
	return c.kitexClient
}

func (c *clientImpl) AIWithOrders(ctx context.Context, Req *AIEino.AIWithOrdersReq, callOptions ...callopt.Option) (r *AIEino.AIWithOrdersResp, err error) {
	return c.kitexClient.AIWithOrders(ctx, Req, callOptions...)
}

func (c *clientImpl) AIWithPreCheckout(ctx context.Context, Req *AIEino.AIWithPreCheckoutReq, callOptions ...callopt.Option) (r *AIEino.AIWithPreCheckoutResp, err error) {
	return c.kitexClient.AIWithPreCheckout(ctx, Req, callOptions...)
}
