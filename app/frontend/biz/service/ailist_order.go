package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gomall/app/frontend/infra/rpc"
	"gomall/app/frontend/utils"
	"gomall/rpc_gen/kitex_gen/AIEino"

	"github.com/cloudwego/hertz/pkg/app"
	common "gomall/app/frontend/hertz_gen/frontend/common"
)

type AIListOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAIListOrderService(Context context.Context, RequestContext *app.RequestContext) *AIListOrderService {
	return &AIListOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *AIListOrderService) Run(req *common.Empty) error {
	userId := utils.GetUserIdFormCtx(h.Context)
	var requestData struct {
		OrderList string `json:"orderList"`
	}
	if err := json.Unmarshal(h.RequestContext.Request.Body(), &requestData); err != nil {
		h.RequestContext.JSON(consts.StatusBadRequest, map[string]string{
			"error": "Invalid JSON data",
		})
		return err
	}
	// 组装请求体
	rpcReq := fmt.Sprintf("需求是: %s; 如果从需求中没有体现出具体的user_id，则使用当前用户的user_id : %d",
		requestData.OrderList, userId)
	rpcResp, err := rpc.AIEinoClient.AIWithOrders(h.Context, &AIEino.AIWithOrdersReq{
		UserInput: rpcReq,
	})
	if err != nil {
		h.RequestContext.JSON(consts.StatusInternalServerError, map[string]string{
			"error": "Internal server error: " + err.Error(),
		})
		return err
	}

	// 确保rpcResp.Orders是可序列化的格式
	h.RequestContext.JSON(consts.StatusOK, map[string]interface{}{
		"orderInfo": rpcResp.Orders,
	})

	return nil
}
