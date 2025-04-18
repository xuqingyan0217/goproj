package localAgent

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/ddgsearch"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/mark3labs/mcp-go/client"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
)

var (
	MCPClient      *client.SSEMCPClient
	mcpClientOnce  sync.Once
	mcpClientMutex sync.Mutex
	err            error
)

// newTool 用于创建 DuckDuckGo 搜索工具组件
func newTool(ctx context.Context) (bt tool.BaseTool, err error) {
	// 配置 DuckDuckGo 工具参数
	config := &duckduckgo.Config{
		MaxResults: 5,                  // 限制返回结果数量为5
		Region:     ddgsearch.RegionCN, // 设置地区为中国
		DDGConfig: &ddgsearch.Config{
			Timeout:    10 * time.Second,                  // 超时时间
			Cache:      false,                             // 不使用缓存
			MaxRetries: 10,                                // 最大重试次数
			Proxy:      os.Getenv("DUCKDUCKGO_PROXY_URL"), // 代理地址
		},
	}
	bt, err = duckduckgo.NewTool(ctx, config)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

/*
/sse 是 MCP 服务器默认的 SSE 端点路径。这个路径是 MCP 框架约定的，用于：
- 区分普通的 HTTP 请求和 SSE 连接
- 标识这是一个 SSE 连接端点
- 保持与 MCP 协议的一致性
*/

// getMCPToolClient 用于创建并初始化 MCP SSE 客户端
func getMCPToolClient(ctx context.Context) (*client.SSEMCPClient, error) {
	mcpClientMutex.Lock()
	defer mcpClientMutex.Unlock()

	mcpClientOnce.Do(func() {
		MCPClient, err = client.NewSSEMCPClient(os.Getenv("MCP_SEE_BASE_URL"))
		if err != nil {
			log.Printf("Failed to create MCP client: %v", err)
			return
		}
		// 启动客户端
		err = MCPClient.Start(ctx)
		if err != nil {
			log.Printf("Failed to start MCP client: %v", err)
			// Clean up the client if start fails
			MCPClient = nil
			return
		}

		// 初始化请求参数
		initRequest := mcpgo.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcpgo.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcpgo.Implementation{
			Name:    "example-client",
			Version: "1.1.0",
		}

		// 发送初始化请求
		_, err = MCPClient.Initialize(ctx, initRequest)
		if err != nil {
			log.Printf("Failed to initialize MCP client: %v", err)
			// Clean up the client if initialize fails
			_ = MCPClient.Close() // Attempt to close the started client
			MCPClient = nil
			return
		}
	})

	if err != nil {
		return nil, fmt.Errorf("MCP client initialization failed: %w", err)
	}
	if MCPClient == nil {
		return nil, fmt.Errorf("MCP client is nil after initialization attempt")
	}

	return MCPClient, nil
}

// CloseMCPClient safely closes the MCP client.
func CloseMCPClient() error {
	mcpClientMutex.Lock()
	defer mcpClientMutex.Unlock()

	if MCPClient != nil {
		err := MCPClient.Close()
		if err != nil {
			log.Printf("Error closing MCP client: %v", err)
			return err
		}
		MCPClient = nil // Set to nil after successful close
		log.Println("MCP client closed successfully")
		return nil
	}
	log.Println("MCP client was already nil or not initialized")
	return nil // No error if already closed or not initialized
}

// Tool1Impl 是 Tool1 的实现结构体
// 可根据实际需求扩展
// config 字段用于存储工具配置
// Tool1Config 用于配置 Tool1
// 目前为空结构体，可根据需求扩展

type Tool1Impl struct {
	config *Tool1Config
}

type Tool1Config struct {
}

// newTool1 用于创建 Tool1 实例
func newTool1(ctx context.Context) (bt tool.BaseTool, err error) {
	// 可在此处修改 Tool1 的配置
	config := &Tool1Config{}
	bt = &Tool1Impl{config: config}
	return bt, nil
}

// Info 返回工具的元信息
func (impl *Tool1Impl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	panic("implement me")
}

// InvokableRun 实现工具的调用逻辑
func (impl *Tool1Impl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	panic("implement me")
}
