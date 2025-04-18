package utils

import (
	"context"
	"gomall/app/frontend/infra/rpc"
	frontentUtils "gomall/app/frontend/utils"
	rpccart "gomall/rpc_gen/kitex_gen/cart"

	"github.com/cloudwego/hertz/pkg/app"
)

// SendErrResponse  pack error response
func SendErrResponse(ctx context.Context, c *app.RequestContext, code int, err error) {
	resp := map[string]interface{}{
		"code":    -1,
		"error":   err.Error(),
		"message": err.Error(),
		"data":    nil,
	}
	c.JSON(code, resp)
}

// SendSuccessResponse  pack success response
func SendSuccessResponse(ctx context.Context, c *app.RequestContext, code int, data interface{}) {
	c.JSON(code, data)
}

func WarpResponse(ctx context.Context, c *app.RequestContext, content map[string]any) map[string]any {
	userId := frontentUtils.GetUserIdFormCtx(ctx)
	content["user_id"] = userId
	if userId > 0 {
		cartResp, err := rpc.CartClient.GetCart(ctx, &rpccart.GetCartReq{
			UserId: uint32(userId),
		})
		if err == nil && cartResp != nil {
			totalNum := uint32(0)
			for _, item := range cartResp.Items {
				totalNum += item.Quantity
			}
			content["cart_num"] = totalNum
		}
	}
	return content
}
