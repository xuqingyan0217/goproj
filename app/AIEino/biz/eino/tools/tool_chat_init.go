package tools

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"log"
)

func LogicAI(ctx context.Context, tools []tool.BaseTool, input string) ([]*schema.Message, error) {
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
	// 第一次调用，获取原始数据
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.System,
			Content: "你是一个能够使用工具的AI助手。当用户提供需求时，你需要理解并使用适当的工具来完成任务。请仔细分析用户的输入，并调用合适的工具来处理。",
		},
		{
			Role:    schema.User,
			Content: input,
		},
	})
	if err != nil {
		log.Fatalf("Failed to invoke agent: %v", err)
	}

	// 第二次调用，对结果进行润色
	if len(resp) > 0 {
		enhancePrompt := "请将以下数据转换成更易读的格式，使用自然语言描述：\n" + resp[0].Content
		enhancedResp, err := chatModel.Generate(ctx, []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个专业的数据分析助手。请将提供的数据转换成更易读的格式，使用自然语言描述，保持专业性和准确性，同时确保信息的完整性。",
			},
			{
				Role:    schema.User,
				Content: enhancePrompt,
			},
		})
		if err != nil {
			log.Printf("Failed to enhance response: %v", err)
			return resp, nil // 如果润色失败，返回原始数据
		}
		return []*schema.Message{enhancedResp}, nil
	}

	return resp, err
}
