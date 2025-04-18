package main

// 实现多个工具之间的调用目前很麻烦，那么到时候就直接用一个工具，外加很多的函数
// 如下实现了一个示例，通过tool，实现ai语句文本内容自动调用rpc商品搜索服务并只返回出商品id
import (
	"context"
	"encoding/json"
	"fmt"
	"gomall/app/AIEino/infra/rpc"
	"gomall/rpc_gen/kitex_gen/product"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type FindOutParams struct {
	Query string `json:"query"`
	// 目前实现的是查询一类商品，每类商品数量都设置为同一个值；因为其他类型不知道为什么不支持
	ProductCount uint32 `json:"product_count"`
}

// ProductInfo 结构体，用于打包返回值
type ProductInfo struct {
	ProductId uint32 `json:"product_id"`
	Quantity  uint32 `json:"quantity"`
}
type FindProductsTool struct{}

func (ft *FindProductsTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "find_products_tool",
		Desc: "Find products by name",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Desc:     "商品名称，作为查询商品的条件，从文本中获取",
				Type:     schema.String,
				Required: true,
			},
			"product_count": {
				Desc:     "商品购买数量，从文本中获取",
				Type:     schema.Integer,
				Required: true,
			},
		}),
	}, nil
}

func (ft *FindProductsTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var params FindOutParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return "", err
	}
	// 验证参数的有效性
	if params.Query == "" {
		return "", fmt.Errorf("name is required")
	}
	// 如果没有从文本中提取到数量，就默认为1件
	if params.ProductCount == 0 {
		params.ProductCount = 1
	}
	// 在此处添加工具的具体逻辑
	var productInfoList []ProductInfo
	productSearchResp, err := rpc.ProductClient.SearchProducts(ctx, &product.SearchProductsReq{
		Query: params.Query,
	})
	if err != nil {
		return "", err
	}
	var productResult = productSearchResp.Results
	for _, pt := range productResult {
		productInfoList = append(productInfoList, ProductInfo{
			ProductId: pt.Id,
			Quantity:  params.ProductCount,
		})
	}

	result, err := json.Marshal(productInfoList)
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
	findProductsTool := &FindProductsTool{}
	tools := []tool.BaseTool{
		findProductsTool,
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
			Content: "使用 find_products_tool 工具查询名称包含 t-shirt 的商品的 ID , 同时每个商品数量设置为 2 ",
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
		if msg.ResponseMeta != nil {
			// 获取 Token 使用情况
			if usage := msg.ResponseMeta.Usage; usage != nil {
				log.Println("提示 Tokens:", usage.PromptTokens)
				log.Println("生成 Tokens:", usage.CompletionTokens)
				log.Println("总 Tokens:", usage.TotalTokens)
			}
		} else {
			log.Println("暂时无法获取token使用量，请前往豆包平台官网查看")
		}
	}
}
