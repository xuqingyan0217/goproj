package main

import (
	"my_assistant/cmd/einoagent/agent"
	"my_assistant/localAgent"
	"my_assistant/pkg/mcp_tools"

	"my_assistant/pkg/env"

	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/devops"

	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"go.opentelemetry.io/otel/attribute"
)

func init() {
	// check some essential envs
	env.MustHasEnvs("ARK_CHAT_MODEL", "ARK_EMBEDDING_MODEL", "ARK_API_KEY")

	if os.Getenv("EINO_DEBUG") != "false" {
		err := devops.Init(context.Background())
		if err != nil {
			log.Printf("[eino dev] init failed, err=%v", err)
		}
	}
}

func main() {
	// 创建一个带取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 获取端口配置
	port := os.Getenv("HERTZ_SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// 启动MCP服务器
	mcpErrChan := make(chan error, 1)
	go func() {
		defer close(mcpErrChan)
		mcp_tools.StartMCPServer()
		select {
		case <-ctx.Done():
			return
		case mcpErrChan <- nil:
		}
	}()

	// 创建 Hertz 服务器
	h := server.Default(server.WithHostPorts(":" + port))

	h.Use(LogMiddleware())

	if os.Getenv("APMPLUS_APP_KEY") != "" {
		region := os.Getenv("APMPLUS_REGION")
		if region == "" {
			region = "cn-beijing"
		}
		_ = provider.NewOpenTelemetryProvider(
			provider.WithServiceName("eino-assistant"),
			provider.WithExportEndpoint(fmt.Sprintf("apmplus-%s.volces.com:4317", region)),
			provider.WithInsecure(),
			provider.WithHeaders(map[string]string{"X-ByteAPM-AppKey": os.Getenv("APMPLUS_APP_KEY")}),
			provider.WithResourceAttribute(attribute.String("apmplus.business_type", "llm")),
		)
		tracer, cfg := hertztracing.NewServerTracer()
		h = server.Default(server.WithHostPorts(":"+port), tracer)
		h.Use(LogMiddleware(), hertztracing.ServerMiddleware(cfg))
	}

	// 注册 agent 路由组
	agentGroup := h.Group("/agent")
	if err := agent.BindRoutes(agentGroup); err != nil {
		log.Fatal("failed to bind agent routes:", err)
	}

	// Redirect root path to /agent
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.Redirect(302, []byte("/agent"))
	})

	// 设置优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待MCP服务器启动完成
	select {
	case err := <-mcpErrChan:
		if err != nil {
			log.Printf("MCP server error: %v", err)
			cancel()
			return
		}
	case <-time.After(3 * time.Second):
		log.Println("MCP server started successfully")
	}

	// 启动主服务器
	h.Spin()
	// 等待信号
	<-sigChan
	log.Println("Shutting down servers...")

	// 优雅关闭
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// 通知所有服务关闭
	cancel()

	// 关闭MCP服务器，包括mcpServer 和 sseServer
	if err := mcp_tools.StopMCPServer(shutdownCtx); err != nil {
		log.Printf("Error during MCP server shutdown: %v", err)
	}
	// 关闭MCP客户端
	if err := localAgent.CloseMCPClient(); err != nil {
		log.Printf("Error during MCP client shutdown: %v", err)
	}
	// 等待Hertz服务器完全关闭
	if err := h.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during Hertz server shutdown: %v", err)
	}
}

// LogMiddleware 记录 HTTP 请求日志
func LogMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Request.URI().Path())
		method := string(c.Request.Method())

		// 处理请求
		c.Next(ctx)

		// 记录请求信息
		latency := time.Since(start)
		statusCode := c.Response.StatusCode()
		log.Printf("[HTTP] %s %s %d %v\n", method, path, statusCode, latency)
	}
}
