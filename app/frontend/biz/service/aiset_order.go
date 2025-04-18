package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	common "gomall/app/frontend/hertz_gen/frontend/common"
	"gomall/app/frontend/infra/rpc"
	"gomall/app/frontend/utils"
	"gomall/rpc_gen/kitex_gen/AIEino"
)

type AISetOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAISetOrderService(Context context.Context, RequestContext *app.RequestContext) *AISetOrderService {
	return &AISetOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *AISetOrderService) Run(req *common.Empty) error {
	// 解析前端发送的JSON数据
	userId := utils.GetUserIdFormCtx(h.Context)
	var requestData struct {
		Order string `json:"order"`
	}
	if err := json.Unmarshal(h.RequestContext.Request.Body(), &requestData); err != nil {
		h.RequestContext.JSON(consts.StatusBadRequest, map[string]string{
			"error": "Invalid JSON data",
		})
		return err
	}
	fmt.Println("req : ", requestData.Order)
	// 组装RPC请求数据
	rpcReq := fmt.Sprintf("使用 pre_checkout_tool 工具; 需求是: %s ,依据需求查找出相应商品id, 返回最终预下单后的结果; 其中基本信息有:下单 user_id 是%d, email 是%s, country 是%s, state 是%s, city 是%s, street_address 是%s, zipcode 是%s",
		requestData.Order, userId, "example@qq.com", "中国", "河南", "南阳", "枣林街", "141588")
	// 调用rpc服务
	rpcResp, err := rpc.AIEinoClient.AIWithPreCheckout(h.Context, &AIEino.AIWithPreCheckoutReq{
		UserInput: rpcReq,
	})
	if err != nil {
		h.RequestContext.JSON(consts.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
		return err
	}

	// 处理RPC响应
	if len(rpcResp.PreCheckoutRes) > 0 {
		// 假设第一个订单号作为返回结果
		orderId := rpcResp.PreCheckoutRes
		fmt.Println(orderId)
		h.RequestContext.JSON(consts.StatusOK, map[string][]string{
			"orderId": orderId,
		})
	} else {
		h.RequestContext.JSON(consts.StatusInternalServerError, map[string]string{
			"error": "No order ID returned from RPC service",
		})
		return errors.New("no order ID returned from RPC service")
	}

	return nil
}
