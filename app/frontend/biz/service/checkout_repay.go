package service

import (
	"context"

	checkout "gomall/app/frontend/hertz_gen/frontend/checkout"
	rpccheckout "gomall/rpc_gen/kitex_gen/checkout"
	rpcpayment "gomall/rpc_gen/kitex_gen/payment"

	"gomall/app/frontend/infra/rpc"
	frontendUtils "gomall/app/frontend/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type CheckoutRepayService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCheckoutRepayService(Context context.Context, RequestContext *app.RequestContext) *CheckoutRepayService {
	return &CheckoutRepayService{RequestContext: RequestContext, Context: Context}
}

func (h *CheckoutRepayService) Run(req *checkout.CheckoutRepayReq) (resp map[string]any, err error) {
	userId := frontendUtils.GetUserIdFormCtx(h.Context)
	repayResp, err := rpc.CheckoutClient.ContinueCheckout(h.Context, &rpccheckout.ContinueCheckoutReq{
		UserId: uint32(userId),
		OrderId: req.OrderId,
		Total: req.Total,
		Email: req.Email,
		CreditCard: &rpcpayment.CreditCardInfo{
			CreditCardNumber:          req.CardNum,
			CreditCardCvv:             req.Cvv,
			CreditCardExpirationYear:  req.ExpirationYear,
			CreditCardExpirationMonth: req.ExpirationMonth,
		},
	})
	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(5005001, err.Error())
	}
	return utils.H{
		"title":    "waiting",
		"redirect": "/checkout/result",
		"transactionId": repayResp.TransactionId,
	}, nil
}
