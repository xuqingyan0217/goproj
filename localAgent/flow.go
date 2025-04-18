package localAgent

import (
	"context"
	tmcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"log"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

// newLambda1 component initialization function of node 'Lambda2' in graph 'localAgent'
func newLambda1(ctx context.Context) (lba *compose.Lambda, err error) {
	// 获取到客户端
	cli, err := getMCPToolClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// 客户端需要在main.go里面在服务关闭的时候一同关闭，提前关闭会影响多次调用工具的执行
	// defer func() { // Removed this block
	// 	err = cli.Close()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()
	// 获取mcp的server端的工具
	mcpTools, err := tmcp.GetTools(ctx, &tmcp.Config{Cli: cli})
	if err != nil {
		log.Fatal(err)
	}

	config := &react.AgentConfig{
		MaxStep:            25,
		ToolReturnDirectly: map[string]struct{}{},
	}
	chatModelIns11, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	config.ToolCallingModel = chatModelIns11
	// 官方tool
	toolIns21, err := newTool(ctx)
	if err != nil {
		return nil, err
	}
	// 一并加入到mcp tool里面
	mcpTools = append(mcpTools, toolIns21)
	config.ToolsConfig.Tools = mcpTools
	ins, err := react.NewAgent(ctx, config)
	if err != nil {
		return nil, err
	}
	lba, err = compose.AnyLambda(ins.Generate, ins.Stream, nil, nil)
	if err != nil {
		return nil, err
	}
	return lba, nil
}
