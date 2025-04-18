package main

// 如下实现了一个示例，通过tool，实现ai语句文本内容自动调用rpc商品搜索服务并只返回出商品id
import (
	"context"
	"encoding/json"
	"fmt"
	"gomall/app/AIEino/infra/rpc"
	"gomall/rpc_gen/kitex_gen/order"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
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

	result, err := json.Marshal(orderResult)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
func main() {
	// 这里因为是单独使用，主服务并没有启动，所以要手动初始化rpc客户端
	rpc.InitClient()
	ctx := context.Background()
	// 初始化工具
	getOrderTool := &GetOrderTool{}
	tools := []tool.BaseTool{
		getOrderTool,
	}

	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: "xxx",
		Model:  "xxx",
	})
	if err != nil {
		log.Fatalf("Failed to create ChatModel: %v", err)
	}

	// 获取工具信息，这个info是间接调用，实际上调用的就是上面实现的方法
	toolInfos := make([]*schema.ToolInfo, 0, len(tools))
	for _, t := range tools {
		info, err := t.Info(ctx)
		if err != nil {
			log.Fatalf("Failed to get tool info: %v", err)
		}
		toolInfos = append(toolInfos, info)
	}

	// 将工具绑定到 ChatModel,绑定时需要用的刚才的工具信息
	err = chatModel.BindTools(toolInfos)
	if err != nil {
		log.Fatalf("Failed to bind tools: %v", err)
	}

	// 创建工具节点
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: tools,
	})
	if err != nil {
		log.Fatalf("Failed to create ToolsNode: %v", err)
	}

	// 构建处理链
	/*graph := compose.NewGraph[[]*schema.Message, []*schema.Message]()
	_ = graph.AddChatModelNode("chat_model", chatModel)
	_ = graph.AddToolsNode("tools", toolsNode)
	*/
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
		AppendToolsNode(toolsNode, compose.WithNodeName("tools"))

	// 编译并运行处理链
	agent, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("Failed to compile chain: %v", err)
	}

	// 运行示例，这里的Invoke其实就是调用的上面的InvokableRun方法
	// 正常输出
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role: schema.User,
			// 这里的文本内容也是有讲究的，最好是带上工具的名称，这样容错高
			Content: "通过 get_order_tool 工具查询 user_id 为 2 的用户的订单列表",
		},
	})
	if err != nil {
		log.Fatalf("Failed to invoke agent: %v", err)
	}
	// 最后是打印结果和Token使用情况，但是不知为何token获取不到
	for _, msg := range resp {
		log.Printf("%s: %s", msg.Role, msg.Content)
		// 获取 Token 使用情况
		// 先检查 ResponseMeta 是否为 nil
		/*if msg.ResponseMeta != nil {
			// 获取 Token 使用情况
			if usage := msg.ResponseMeta.Usage; usage != nil {
				log.Println("提示 Tokens:", usage.PromptTokens)
				log.Println("生成 Tokens:", usage.CompletionTokens)
				log.Println("总 Tokens:", usage.TotalTokens)
			}
		} else {
			log.Println("暂时无法获取token使用量，请前往豆包平台官网查看")
		}*/
	}
}
