package service

import (
	"context"
	"fmt"

	"gomall/app/checkout/infra/mq"
	"gomall/app/checkout/infra/rpc"
	"gomall/rpc_gen/kitex_gen/cart"
	checkout "gomall/rpc_gen/kitex_gen/checkout"
	"gomall/rpc_gen/kitex_gen/email"
	"gomall/rpc_gen/kitex_gen/order"
	"gomall/rpc_gen/kitex_gen/payment"
	"gomall/rpc_gen/kitex_gen/product"
	"strconv"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
)

type CheckoutService struct {
	ctx context.Context
} // NewCheckoutService new CheckoutService
func NewCheckoutService(ctx context.Context) *CheckoutService {
	return &CheckoutService{ctx: ctx}
}

// Run create note info
func (s *CheckoutService) Run(req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	// Finish your business logic.
	// 获取购物车信息
	cartResult, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(5005001, err.Error())
	}
	if cartResult == nil || cartResult.Items == nil {
		return nil, kerrors.NewGRPCBizStatusError(5004001, "cart is empty")
	}

	// 计算总额
	var (
		total float32
		oi    []*order.OrderItem
	)
	for _, cartItem := range cartResult.Items {
		productResp, resultErr := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{
			Id: cartItem.ProductId,
		})
		if resultErr != nil {
			return nil, kerrors.NewGRPCBizStatusError(5005002, resultErr.Error())
		}
		if productResp.Product == nil {
			continue
		}
		p := productResp.Product.Price
		cost := p * float32(cartItem.Quantity)
		total += cost
		// 订单信息相关
		oi = append(oi, &order.OrderItem{
			Item: &cart.CartItem{
				ProductId: cartItem.ProductId,
				Quantity:  cartItem.Quantity,
			},
			Cost: cost,
		})
	}

	// 创建订单
	var orderId string
	// 使用 strconv.ParseInt 进行类型转换
	int64ZipCode, err := strconv.ParseInt(req.Address.ZipCode, 10, 32)
	if err != nil {
		return
	}
	zipcode := int32(int64ZipCode)
	orderResp, err := rpc.OrderClient.PlaceOrder(s.ctx, &order.PlaceOrderReq{
		UserId: req.UserId,
		Email:  req.Email,
		Address: &order.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       zipcode,
		},
		Items: oi,
		UserCurrency: "ORDERED",
	})
	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(5005005, err.Error())
	}
	if orderResp == nil || orderResp.Order == nil {
		return nil, kerrors.NewGRPCBizStatusError(5005006, "order is empty")
	}
	orderId = orderResp.Order.OrderId

	// 如果是取消支付，则直接返回订单信息
	if req.Flag == 0 {
		// 清空购物车
		_, err = rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{
			UserId: req.UserId,
		})
		if err != nil {
			return nil, kerrors.NewGRPCBizStatusError(5005003, err.Error())
		}
		status, err := rpc.OrderClient.CancelPayment(s.ctx, &order.CancelPaymentReq{
			OrderId: orderId,
			Status: "WAIT",
		})
		if err != nil {
			return nil, kerrors.NewGRPCBizStatusError(5005007, err.Error())
		}
		if status.Success {
			return nil, kerrors.NewGRPCBizStatusError(5005008, "cancel payment")
		}
	}

	// 支付校验参数
	payReq := &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: orderId,
		Amount:  total,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
		},
	}

	// 清空购物车
	_, err = rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(5005003, err.Error())
	}
	fmt.Println("--------------------------010-----------------------")
	// 调用支付
	paymentResult, err := rpc.PaymentClient.Charge(s.ctx, payReq)
	if err != nil {
		fmt.Println("--------------------------err-----------------------", err)
		return nil, kerrors.NewGRPCBizStatusError(5005004, err.Error())
	}
	klog.Info("paymentResult:", paymentResult)

	// 下单后发送邮件
	// 生产消息
	data, _ := proto.Marshal(&email.EmailReq{
		From:        "from@email.com",
		To:          req.Email,
		ContentType: "text/plain",
		Subject:     "You hava just created an order in the XuQy shop",
		Content:     "You hava just created an order in the XuQy shop",
	})
	msg := &nats.Msg{
		Subject: "email",
		Data:    data,
		Header:  make(nats.Header),
	}
	otel.GetTextMapPropagator().Inject(s.ctx, propagation.HeaderCarrier(msg.Header))
	_ = mq.Nc.PublishMsg(msg)

	// 返回相应
	resp = &checkout.CheckoutResp{
		OrderId:       orderId,
		TransactionId: paymentResult.TransactionId,
	}
	return
}
