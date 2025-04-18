package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	einoagent "my_assistant/localAgent"

	"io"
	"log"
	"os"
	"sync"

	"github.com/cloudwego/eino-ext/callbacks/apmplus"
	"github.com/cloudwego/eino-ext/callbacks/langfuse"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	"my_assistant/pkg/mem"
)

// memory 变量用于存储默认内存实例，以便在全局范围内使用。
var memory = mem.GetDefaultMemory()

// cbHandler 变量用于处理回调事件，确保回调逻辑与业务代码分离。
var cbHandler callbacks.Handler

// once 变量用于确保某些初始化操作只被执行一次，避免并发环境下的重复执行。
var once sync.Once

// Init 初始化日志和全局回调处理器。
// 该函数确保初始化逻辑仅执行一次，避免重复初始化。
// 主要功能包括：
//   - 创建日志目录并初始化日志文件。
//   - 配置日志回调处理器，支持调试模式。
//   - 根据环境变量初始化 APMPlus 和 Langfuse 回调处理器。
//   - 注册回调处理器以支持跟踪和指标功能。
//
// 返回值:
//   - error: 如果初始化过程中发生错误，则返回具体的错误信息。
func Init() error {
	var err error
	once.Do(func() {
		// 确保日志目录存在，如果不存在则创建。
		os.MkdirAll("log", 0755)
		var f *os.File
		// 打开或创建日志文件，用于写入日志数据。
		f, err = os.OpenFile("log/eino.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return
		}

		// 配置日志回调的详细设置。
		cbConfig := &LogCallbackConfig{
			Detail: true,
			Writer: f,
		}
		// 如果 DEBUG 环境变量为 "true"，启用调试模式。
		if os.Getenv("DEBUG") == "true" {
			cbConfig.Debug = true
		}
		// 初始化日志回调处理器。
		cbHandler = LogCallback(cbConfig)

		// 初始化全局回调处理器列表，用于跟踪和指标功能。
		callbackHandlers := make([]callbacks.Handler, 0)
		// 如果设置了 APMPlus 应用密钥，则初始化 APMPlus 回调处理器。
		if os.Getenv("APMPLUS_APP_KEY") != "" {
			region := os.Getenv("APMPLUS_REGION")
			if region == "" {
				region = "cn-beijing"
			}
			fmt.Println("[eino agent] INFO: use apmplus as callback, watch at: https://console.volcengine.com/apmplus-server")
			cbh, _, err := apmplus.NewApmplusHandler(&apmplus.Config{
				Host:        fmt.Sprintf("apmplus-%s.volces.com:4317", region),
				AppKey:      os.Getenv("APMPLUS_APP_KEY"),
				ServiceName: "eino-assistant",
				Release:     "release/v0.0.1",
			})
			if err != nil {
				log.Fatal(err)
			}

			callbackHandlers = append(callbackHandlers, cbh)
		}

		// 如果设置了 Langfuse 公钥和私钥，则初始化 Langfuse 回调处理器。
		if os.Getenv("LANGFUSE_PUBLIC_KEY") != "" && os.Getenv("LANGFUSE_SECRET_KEY") != "" {
			fmt.Println("[eino agent] INFO: use langfuse as callback, watch at: https://cloud.langfuse.com")
			cbh, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
				Host:      "https://cloud.langfuse.com",
				PublicKey: os.Getenv("LANGFUSE_PUBLIC_KEY"),
				SecretKey: os.Getenv("LANGFUSE_SECRET_KEY"),
				Name:      "Eino Assistant",
				Public:    true,
				Release:   "release/v0.0.1",
				UserID:    "eino_god",
				Tags:      []string{"eino", "assistant"},
			})
			callbackHandlers = append(callbackHandlers, cbh)
		}
		// 如果存在任何回调处理器，则初始化全局回调处理器。
		if len(callbackHandlers) > 0 {
			// InitCallbackHandlers 被遗弃了
			// callbacks.InitCallbackHandlers(callbackHandlers)
			callbacks.AppendGlobalHandlers(callbackHandlers...)
		}
	})
	return err
}

// RunAgent 初始化并运行一个代理来处理用户消息，支持流式响应。
// 该函数构建代理、管理对话历史，并将对话保存到内存中。
// 返回一个用于接收代理响应的流读取器和可能发生的错误。
//
// 参数:
// - ctx: 上下文，用于控制取消和超时。
// - id: 对话或会话的唯一标识符。
// - msg: 用户输入的消息，由代理进行处理。
//
// 返回值:
// - *schema.StreamReader[*schema.Message]: 用于接收代理响应的流读取器。
// - error: 如果代理初始化或流式传输失败，则返回错误。
func RunAgent(ctx context.Context, id string, msg string) (*schema.StreamReader[*schema.Message], error) {

	// 使用上下文构建代理运行器。
	runner, err := einoagent.BuildLocalAgent(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build agent graph: %w", err)
	}

	// 获取或创建与指定 ID 关联的对话对象。
	conversation := memory.GetConversation(id, true)

	// 构造包含用户查询和对话历史的 UserMessage 对象。
	userMessage := &einoagent.UserMessage{
		ID:      id,
		Query:   msg,
		History: conversation.GetMessages(),
	}

	// 使用回调处理器通过代理运行器流式传输用户消息。
	sr, err := runner.Stream(ctx, userMessage, compose.WithCallbacks(cbHandler))
	if err != nil {
		return nil, fmt.Errorf("failed to stream: %w", err)
	}

	// 创建流读取器的两个副本：一个用于返回给调用者，另一个用于内部处理（保存到内存）。
	srs := sr.Copy(2)

	// 启动 Goroutine 处理流式数据并保存对话到内存。
	go func() {
		// 收集所有流式消息块以便后续合并。
		fullMsgs := make([]*schema.Message, 0)

		// 最后把整合后的消息保存起来
		defer func() {
			// 关闭内部使用的流读取器。
			srs[1].Close()

			// 将用户输入添加到对话历史中。
			conversation.Append(schema.UserMessage(msg))

			// 合并所有接收到的消息块为一条完整消息。
			fullMsg, err := schema.ConcatMessages(fullMsgs)
			if err != nil {
				fmt.Println("error concatenating messages: ", err.Error())
			}

			// 将代理的完整响应添加到对话历史中。
			conversation.Append(fullMsg)
		}()

		// 循环处理从流中接收到的消息块。
	outer:
		for {
			select {
			case <-ctx.Done():
				// 如果上下文被取消或超时，退出循环。
				fmt.Println("context done", ctx.Err())
				return
			default:
				// 从流中接收下一个消息块。
				chunk, err := srs[1].Recv()
				if err != nil {
					// 如果流结束（EOF）或发生错误，退出循环。
					if errors.Is(err, io.EOF) {
						break outer
					}
				}

				// 将接收到的消息块添加到完整消息列表中。
				fullMsgs = append(fullMsgs, chunk)
			}
		}
	}()

	// 返回主流读取器供调用者使用。
	return srs[0], nil
}

type LogCallbackConfig struct {
	Detail bool
	Debug  bool
	Writer io.Writer
}

// LogCallback 创建并返回一个日志回调处理程序，用于在执行前后记录组件的信息。
// 它接受一个指向 LogCallbackConfig 结构的指针作为参数，以配置日志的详细程度和输出目标。
// 如果未提供配置或配置为nil，则使用默认配置，其中详细模式为true，输出目标为标准输出。
func LogCallback(config *LogCallbackConfig) callbacks.Handler {
	// 检查是否提供了配置，如果没有，则初始化配置为默认值。
	if config == nil {
		config = &LogCallbackConfig{
			Detail: true,
			Writer: os.Stdout,
		}
	}
	// 确保配置中的Writer不为nil，如果为nil，则设置为标准输出。
	if config.Writer == nil {
		config.Writer = os.Stdout
	}
	// 创建一个新的回调处理程序构建器。
	builder := callbacks.NewHandlerBuilder()
	// 设置处理程序在执行开始时的动作。
	builder.OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		// 记录开始执行的信息，包括组件、类型和名称。
		fmt.Fprintf(config.Writer, "[view]: start [%s:%s:%s]\n", info.Component, info.Type, info.Name)
		// 如果配置了详细模式，則根据配置的Debug选项以适当格式输出输入信息。
		if config.Detail {
			var b []byte
			// 如果启用了调试模式，则使用缩进格式序列化输入信息，以提高可读性。
			if config.Debug {
				b, _ = json.MarshalIndent(input, "", "  ")
			} else {
				// 否则，直接序列化输入信息。
				b, _ = json.Marshal(input)
			}
			// 输出序列化后的输入信息。
			fmt.Fprintf(config.Writer, "%s\n", string(b))
		}
		return ctx
	})
	// 设置处理程序在执行结束时的动作。
	builder.OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		// 记录执行结束的信息，包括组件、类型和名称。
		fmt.Fprintf(config.Writer, "[view]: end [%s:%s:%s]\n", info.Component, info.Type, info.Name)
		return ctx
	})
	// 构建并返回配置好的回调处理程序。
	return builder.Build()
}
