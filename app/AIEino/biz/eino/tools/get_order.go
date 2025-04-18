package tools

// 如下实现了一个示例，通过tool，实现ai语句文本内容自动调用rpc商品搜索服务并只返回出商品id
import (
	"context"
	"encoding/json"
	"fmt"
	"gomall/app/AIEino/infra/rpc"
	"gomall/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/eino/components/tool"

	"github.com/cloudwego/eino/schema"
)

type GetOrderParams struct {
	UserId uint32 `json:"user_id"`
}
type GetOrderTool struct{}

// Info 搭配上述参数添加
func (g *GetOrderTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_order_tool",
		Desc: "Get order by user_id",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"user_id": {
				Desc:     "用户id，用于查询用户的订单列表",
				Type:     schema.Integer,
				Required: true,
			},
		}),
	}, nil
}

// InvokableRun 作为业务核心，当参数进入这里调用rpc接口后，返回值就已经是我们需要的结构了，直接转成string返回即可
func (g *GetOrderTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var params GetOrderParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return "", err
	}
	// 验证参数的有效性
	if params.UserId == 0 {
		return "", fmt.Errorf("user_id is required")
	}

	orderResp, err := rpc.OrderClient.ListOrder(ctx, &order.ListOrderReq{
		UserId: params.UserId,
	})
	if err != nil {
		return "", err
	}
	var orderResult = orderResp.Orders
	fmt.Println(orderResult)

	result, err := json.Marshal(orderResult)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
