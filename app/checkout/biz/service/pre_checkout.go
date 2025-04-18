package service

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"gomall/app/checkout/infra/mq"
	"gomall/rpc_gen/kitex_gen/email"
	"google.golang.org/protobuf/proto"
	"log"
	"time"

	"gomall/rpc_gen/kitex_gen/cart"
	checkout "gomall/rpc_gen/kitex_gen/checkout"
	"gomall/rpc_gen/kitex_gen/product"

	"gomall/app/checkout/infra/rpc"

	"gomall/rpc_gen/kitex_gen/order"
	"strconv"
)

type PreCheckoutService struct {
	ctx context.Context
} // NewPreCheckoutService new PreCheckoutService
func NewPreCheckoutService(ctx context.Context) *PreCheckoutService {
	return &PreCheckoutService{ctx: ctx}
}

// Run create note info
func (s *PreCheckoutService) Run(req *checkout.PreCheckoutReq) (resp *checkout.PreCheckoutResp, err error) {
	// Finish your business logic.
	for _, productItem := range req.ProductInfoList {
		_, err := rpc.CartClient.AddItem(s.ctx, &cart.AddItemReq{
			UserId: req.UserId,
			Item: &cart.CartItem{
				ProductId: productItem.ProductId,
				Quantity:  productItem.Quantity,
			},
		})
		if err != nil {
			return nil, kerrors.NewGRPCBizStatusError(5005003, err.Error())
		}
	}

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

	// 创建预订单
	var orderId string
	// 使用 strconv.ParseInt 进行类型转换
	int64ZipCode, err := strconv.ParseInt(req.Address.ZipCode, 10, 32)
	if err != nil {
		return
	}
	zipcode := int32(int64ZipCode)
	log.Printf("zipcode: %d", zipcode)
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
		UserCurrency: "WAIT",
	})
	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(5005005, err.Error())
	}
	if orderResp == nil || orderResp.Order == nil {
		return nil, kerrors.NewGRPCBizStatusError(5005006, "pre-order is empty")
	}
	orderId = orderResp.Order.OrderId

	// 不进行支付操作，移除支付相关代码

	// 清空购物车
	_, err = rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(5005003, err.Error())
	}
	// 预下单后发送邮件通知用户预订单已生成
	// 生产消息
	data, _ := proto.Marshal(&email.EmailReq{
		From:        "from@email.com",
		To:          req.Email,
		ContentType: "text/plain",
		Subject:     "You have just created a pre-order in the XuQy shop",
		Content:     "You have just created a pre-order in the XuQy shop. Please complete the payment in time.",
	})
	msg := &nats.Msg{
		Subject: "email",
		Data:    data,
		Header:  make(nats.Header),
	}
	otel.GetTextMapPropagator().Inject(s.ctx, propagation.HeaderCarrier(msg.Header))
	_ = mq.Nc.PublishMsg(msg)

	// 返回相应
	resp = &checkout.PreCheckoutResp{
		PreOrderId: orderId,
		//TransactionId: "", // 预下单没有交易ID
		TotalAmount: total, // 返回预订单总额
		ValidUntil:  time.Now().Unix() + 600,
	}
	return
}
