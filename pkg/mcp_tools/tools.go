package mcp_tools

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mark3labs/mcp-go/mcp"
)

var (
	mcpServer   *server.MCPServer
	sseServer   *server.SSEServer
	serverMutex sync.Mutex
)

// StartMCPServer 启动MCP SSE服务，注册工具和Prompt模板
func StartMCPServer() {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	mcpServer = server.NewMCPServer("demo", mcpgo.LATEST_PROTOCOL_VERSION)

	// 自定义回调函数类型
	type ToolCallback func(ctx context.Context, request mcpgo.CallToolRequest, result interface{})

	// 定义一个全局的回调函数
	var toolCallback ToolCallback = func(ctx context.Context, request mcpgo.CallToolRequest, result interface{}) {
		log.Printf("Callback executed for tool %s with result: %v", request.Params.Name, result)
	}

	// 添加新的工具：计算
	mcpServer.AddTool(mcpgo.NewTool("calculate",
		mcpgo.WithDescription("Perform basic arithmetic operations"),
		mcpgo.WithString("operation",
			mcpgo.Required(),
			mcpgo.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcpgo.Enum("add", "subtract", "multiply", "divide"),
		),
		mcpgo.WithNumber("x",
			mcpgo.Required(),
			mcpgo.Description("First number"),
		),
		mcpgo.WithNumber("y",
			mcpgo.Required(),
			mcpgo.Description("Second number"),
		),
	), func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)
		x := request.Params.Arguments["x"].(float64)
		y := request.Params.Arguments["y"].(float64)

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			if y == 0 {
				return mcpgo.NewToolResultError("Cannot divide by zero"), nil
			}
			result = x / y
		}
		toolCallback(ctx, request, result)
		log.Printf("Calculated result: %.2f", result)
		return mcpgo.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})

	// 新增工具：生成两个随机数
	mcpServer.AddTool(mcpgo.NewTool("generateRandomNumbers",
		mcpgo.WithDescription("Generate two random numbers"),
	), func(ctx context.Context, request mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		num1 := r.Float64()
		num2 := r.Float64()
		result := fmt.Sprintf("Random numbers: %.2f, %.2f", num1, num2)
		toolCallback(ctx, request, result)
		return mcpgo.NewToolResultText(result), nil
	})

	// 添加新的工具：连接数据库并获取表结构信息
	mcpServer.AddTool(mcp.NewTool("getDatabaseTableStructure",
		mcp.WithDescription("Connect to a database and get all table structures and information"),
		mcp.WithString("dsn",
			mcp.Required(),
			mcp.Description("Data Source Name for database connection, e.g. user:password@tcp(127.0.0.1:3306)/dbname"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dsn := request.Params.Arguments["dsn"].(string)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to connect to database: %v", err)), nil
		}
		defer db.Close()

		// 检查数据库连接是否正常
		err = db.PingContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to ping database: %v", err)), nil
		}

		// 获取所有表名
		rows, err := db.QueryContext(ctx, "SHOW TABLES")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get tables: %v", err)), nil
		}
		defer rows.Close()

		var tableNames []string
		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to scan table name: %v", err)), nil
			}
			tableNames = append(tableNames, tableName)
		}

		var tableStructures string
		for _, tableName := range tableNames {
			// 获取表结构信息
			tableRows, err := db.QueryContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s", tableName))
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get table structure for %s: %v", tableName, err)), nil
			}

			if tableRows.Next() {
				var createTableSQL string
				var tempTableName string
				if err := tableRows.Scan(&tempTableName, &createTableSQL); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Failed to scan table structure for %s: %v", tableName, err)), nil
				}
				tableStructures += fmt.Sprintf("Table: %s\n%s\n\n", tableName, createTableSQL)
			}
			tableRows.Close()
		}

		return mcp.NewToolResultText(tableStructures), nil
	})

	// 启动 SSE 服务器
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println(e)
			}
		}()
		sseServer = server.NewSSEServer(mcpServer, server.WithBaseURL(os.Getenv("MCP_SSE_SERVER")))
		err := sseServer.Start(os.Getenv("MCP_SSE_ADDR"))
		if err != nil {
			log.Fatal(err)
		}
	}()
}

// StopMCPServer 优雅关闭MCP服务器
func StopMCPServer(ctx context.Context) error {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if sseServer != nil {
		if err := sseServer.Shutdown(ctx); err != nil {
			log.Printf("Error stopping SSE server: %v", err)
			return err
		}
		sseServer = nil
	}

	if mcpServer != nil {
		mcpServer = nil
	}

	return nil
}
