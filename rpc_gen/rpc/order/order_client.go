package order

import (
	"context"
	order "gomall/rpc_gen/kitex_gen/order"

	"gomall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() orderservice.Client
	Service() string
	PlaceOrder(ctx context.Context, Req *order.PlaceOrderReq, callOptions ...callopt.Option) (r *order.PlaceOrderResp, err error)
	ListOrder(ctx context.Context, Req *order.ListOrderReq, callOptions ...callopt.Option) (r *order.ListOrderResp, err error)
	CancelOrder(ctx context.Context, Req *order.CancelOrderReq, callOptions ...callopt.Option) (r *order.CancelOrderResp, err error)
	CancelPayment(ctx context.Context, Req *order.CancelPaymentReq, callOptions ...callopt.Option) (r *order.CancelPaymentResp, err error)
	ChangeOrderStatus(ctx context.Context, Req *order.ChangeOrderStatusReq, callOptions ...callopt.Option) (r *order.ChangeOrderStatusResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := orderservice.NewClient(dstService, opts...)
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
	kitexClient orderservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() orderservice.Client {
	return c.kitexClient
}

func (c *clientImpl) PlaceOrder(ctx context.Context, Req *order.PlaceOrderReq, callOptions ...callopt.Option) (r *order.PlaceOrderResp, err error) {
	return c.kitexClient.PlaceOrder(ctx, Req, callOptions...)
}

func (c *clientImpl) ListOrder(ctx context.Context, Req *order.ListOrderReq, callOptions ...callopt.Option) (r *order.ListOrderResp, err error) {
	return c.kitexClient.ListOrder(ctx, Req, callOptions...)
}

func (c *clientImpl) CancelOrder(ctx context.Context, Req *order.CancelOrderReq, callOptions ...callopt.Option) (r *order.CancelOrderResp, err error) {
	return c.kitexClient.CancelOrder(ctx, Req, callOptions...)
}

func (c *clientImpl) CancelPayment(ctx context.Context, Req *order.CancelPaymentReq, callOptions ...callopt.Option) (r *order.CancelPaymentResp, err error) {
	return c.kitexClient.CancelPayment(ctx, Req, callOptions...)
}

func (c *clientImpl) ChangeOrderStatus(ctx context.Context, Req *order.ChangeOrderStatusReq, callOptions ...callopt.Option) (r *order.ChangeOrderStatusResp, err error) {
	return c.kitexClient.ChangeOrderStatus(ctx, Req, callOptions...)
}
