package agent

import (
	"bufio"

	"my_assistant/pkg/mem"

	"context"
	"embed"
	"errors"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/hertz-contrib/sse"
)

//go:embed web
var webContent embed.FS

type ChatRequest struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func BindRoutes(r *route.RouterGroup) error {
	if err := Init(); err != nil {
		return err
	}

	// API 路由
	r.GET("/api/chat", HandleChat)
	r.GET("/api/log", HandleLog)
	r.GET("/api/history", HandleHistory)
	r.DELETE("/api/history", HandleDeleteHistory)

	// 静态文件服务
	r.GET("/", func(ctx context.Context, c *app.RequestContext) {
		content, err := webContent.ReadFile("web/index.html")
		if err != nil {
			c.String(consts.StatusNotFound, "File not found")
			return
		}
		c.Header("Content-Type", "text/html")
		c.Write(content)
	})

	// 处理GET请求，尝试读取和返回指定的文件内容。
	/*
		这是一个处理静态文件请求的路由处理器。当用户访问类似 /style.css 或 /script.js 这样的URL时，这个处理器会从嵌入的web目录中读取对应的文件内容。它首先通过c.Param("file")获取URL中的文件名，然后尝试从web目录读取该文件。如果文件存在，它会根据文件扩展名设置正确的Content-Type响应头，并将文件内容发送给客户端；如果文件不存在，则返回404错误。
	*/
	// 该路由接受一个文件名参数，并根据该参数读取文件内容，然后将内容作为HTTP响应体返回。
	// 如果文件不存在，将返回404状态码和错误信息。
	r.GET("/:file", func(ctx context.Context, c *app.RequestContext) {
		// 获取URL参数中的文件名。
		file := c.Param("file")

		// 读取指定文件的内容。
		content, err := webContent.ReadFile("web/" + file)
		if err != nil {
			// 如果文件读取失败（例如文件不存在），返回404错误码和错误信息。
			c.String(consts.StatusNotFound, "File not found")
			return
		}

		// 根据文件扩展名确定Content-Type。
		contentType := mime.TypeByExtension(filepath.Ext(file))
		if contentType == "" {
			// 如果无法确定Content-Type，使用默认值。
			contentType = "application/octet-stream"
		}
		// 设置响应头中的Content-Type。
		c.Header("Content-Type", contentType)

		// 将文件内容写入响应体。
		c.Write(content)
	})

	return nil
}

func HandleChat(ctx context.Context, c *app.RequestContext) {
	// 获取请求参数
	id := c.Query("id")
	message := c.Query("message")

	// 参数校验
	if id == "" || message == "" {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"status": "error",
			"error":  "missing id or message parameter",
		})
		return
	}

	log.Printf("[Chat] Starting chat with ID: %s, Message: %s\n", id, message)

	// 启动 AI 代理处理会话
	sr, err := RunAgent(ctx, id, message)
	if err != nil {
		log.Printf("[Chat] Error running agent: %v\n", err)
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// 创建 SSE 事件流 (Server-Sent Events)
	s := sse.NewStream(c)
	defer func() {
		sr.Close() // 关闭流式读取器
		c.Flush()  // 强制刷新响应缓冲区
		log.Printf("[Chat] Finished chat with ID: %s\n", id)
	}()

outer:
	for {
		select {
		case <-ctx.Done(): // 处理上下文取消信号
			log.Printf("[Chat] Context done for chat ID: %s\n", id)
			return
		default:
			// 接收 AI 代理的流式响应
			msg, err := sr.Recv()
			if errors.Is(err, io.EOF) { // 流结束信号
				log.Printf("[Chat] EOF received for chat ID: %s\n", id)
				break outer
			}
			if err != nil { // 处理接收错误
				log.Printf("[Chat] Error receiving message: %v\n", err)
				break outer
			}

			// 通过 SSE 向客户端推送消息
			err = s.Publish(&sse.Event{
				Data: []byte(msg.Content), // 仅发送消息内容
			})
			if err != nil { // 处理推送错误
				log.Printf("[Chat] Error publishing message: %v\n", err)
				break outer
			}
		}
	}
}

func HandleHistory(ctx context.Context, c *app.RequestContext) {
	// query: id => get history, none => list all
	id := c.Query("id")

	if id == "" {
		ids := mem.GetDefaultMemory().ListConversations()

		c.JSON(consts.StatusOK, map[string]interface{}{
			"ids": ids,
		})
		return
	}

	conversation := mem.GetDefaultMemory().GetConversation(id, false)
	if conversation == nil {
		c.JSON(consts.StatusNotFound, map[string]string{
			"error": "conversation not found",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"conversation": conversation,
	})

}

func HandleDeleteHistory(ctx context.Context, c *app.RequestContext) {
	id := c.Query("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"error": "missing id parameter",
		})
		return
	}

	mem.GetDefaultMemory().DeleteConversation(id)
	c.JSON(consts.StatusOK, map[string]string{
		"status": "success",
	})
}

// HandleLog 处理日志流请求。
// 该函数实时读取日志文件，并通过服务器发送事件（SSE）将每一行日志推送到客户端。
// 参数：
//
//	ctx - 请求的上下文，用于在需要时取消操作。
//	c - 请求上下文，包含处理 HTTP 请求和响应所需的信息。
func HandleLog(ctx context.Context, c *app.RequestContext) {
	// 尝试打开日志文件
	file, err := os.Open("log/eino.log")
	if err != nil {
		// 如果打开文件失败，返回错误响应
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	// 确保函数退出时关闭文件
	defer file.Close()

	// 创建一个新的 SSE 流
	s := sse.NewStream(c)
	// 确保函数退出时刷新 HTTP 响应
	defer c.Flush()

	// 定位到文件末尾
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		// 如果定位失败，记录错误并退出函数
		log.Println("定位文件时出错:", err)
		return
	}

	// 使用 Goroutine 持续读取新行
	go func() {
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				// 如果读取失败且不是 EOF，记录错误并跳出循环
				log.Println("读取日志时出错:", err)
				break
			}

			// 如果读取到一行数据，发布它
			if line != "" {
				err = s.Publish(&sse.Event{
					Data: []byte(line),
				})
				if err != nil {
					// 如果发布失败，记录错误并跳出循环
					log.Println("发布日志时出错:", err)
					break
				}
			}

			// 如果到达 EOF，稍作等待后重试
			if err == io.EOF {
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}
	}()

	// 保持连接开启
	<-ctx.Done()
}
