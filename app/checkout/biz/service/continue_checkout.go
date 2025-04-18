package service

import (
	"context"
	"fmt"
	checkout "gomall/rpc_gen/kitex_gen/checkout"
	order "gomall/rpc_gen/kitex_gen/order"
	"gomall/rpc_gen/kitex_gen/payment"
	"gomall/rpc_gen/kitex_gen/email"
	"gomall/app/checkout/infra/mq"
	"gomall/app/checkout/infra/rpc"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
)

type ContinueCheckoutService struct {
	ctx context.Context
} // NewContinueCheckoutService new ContinueCheckoutService
func NewContinueCheckoutService(ctx context.Context) *ContinueCheckoutService {
	return &ContinueCheckoutService{ctx: ctx}
}

// Run create note info
func (s *ContinueCheckoutService) Run(req *checkout.ContinueCheckoutReq) (resp *checkout.ContinueCheckoutResp, err error) {
	// Finish your business logic.
	// 支付校验参数
	payReq := &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
		Amount:  req.Total,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
		},
	}
	// 调用支付
	paymentResult, err := rpc.PaymentClient.Charge(s.ctx, payReq)
	if err != nil {
		fmt.Println("--------------------------err-----------------------", err)
		return nil, kerrors.NewGRPCBizStatusError(5005004, err.Error())
	}
	klog.Info("paymentResult:", paymentResult)
	// 修改order的user_currency字段
	orderStatus, err := rpc.OrderClient.ChangeOrderStatus(s.ctx, &order.ChangeOrderStatusReq{
		OrderId: req.OrderId,
		Status:  "ORDERED",
	})
	if err!= nil {
		return nil, kerrors.NewGRPCBizStatusError(5005005, err.Error())
	}
	if !orderStatus.Success {
		return nil, kerrors.NewBizStatusError(5005006, "change order status failed")
	}
	// 下单后发送邮件
	// 生产消息
	data, _ := proto.Marshal(&email.EmailReq{
		From:        "from@email.com",
		To:          req.Email,
		ContentType: "text/plain",
		Subject:     "You hava Reapayed in the XuQy shop",
		Content:     "You hava Reapayed in the XuQy shop",
	})
	msg := &nats.Msg{
		Subject: "email",
		Data:    data,
		Header:  make(nats.Header),
	}
	otel.GetTextMapPropagator().Inject(s.ctx, propagation.HeaderCarrier(msg.Header))
	_ = mq.Nc.PublishMsg(msg)

	// 返回相应
	return &checkout.ContinueCheckoutResp{
		TransactionId: paymentResult.TransactionId,
	}, nil
}
