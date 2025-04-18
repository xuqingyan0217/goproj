package order

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gomall/app/frontend/biz/service"
	"gomall/app/frontend/biz/utils"
	common "gomall/app/frontend/hertz_gen/frontend/common"
	order "gomall/app/frontend/hertz_gen/frontend/order"
)

// OrderList .
// @router /order [GET]
func OrderList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req common.Empty
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewOrderListService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	// c.HTML(consts.StatusOK, "order", utils.WarpResponse(ctx, c, resp))
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, utils.WarpResponse(ctx, c, resp))
}

// CancelOrder .
// @router /order/cancel [POST]
func CancelOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req order.CancelOrderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewCancelOrderService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, utils.WarpResponse(ctx, c, resp))
}
