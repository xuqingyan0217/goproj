package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"

	"my_assistant/locationKnowledge"
	"my_assistant/pkg/env"

	"github.com/cloudwego/eino/components/embedding"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/document"
)

func init() {
	// check some essential envs
	env.MustHasEnvs("ARK_API_KEY", "ARK_EMBEDDING_MODEL")
}

func main() {
	ctx := context.Background()

	err := indexFiles(ctx, os.Getenv("KNOWLEDGE_BASE_PATH"))
	if err != nil {
		panic(err)
	}

	fmt.Println("index success")
}

// indexMarkdownFiles 索引给定目录下的所有 Markdown 文件。
// 该函数首先构建知识索引结构，然后遍历指定目录中的所有文件。
// 对于每个 Markdown 文件，它将文件内容传递给构建的索引结构进行处理。
// ctx: 上下文，用于取消操作。
// dir: 要索引的目录路径。
// 返回可能发生的错误。
func indexFiles(ctx context.Context, dir string) error {
	// 构建知识索引结构。
	runner, err := locationKnowledge.BuildLocationKnowledge(ctx)
	if err != nil {
		return fmt.Errorf("build index graph failed: %w", err)
	}

	// 遍历 dir 下的所有文件
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk dir failed: %w", err)
		}
		if d.IsDir() {
			return nil
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".md" && ext != ".xlsx" {
			fmt.Printf("[skip] unsupported file type: %s\n", path)
			return nil
		}

		fmt.Printf("[start] indexing file: %s\n", path)

		// 使用构建的索引结构处理文件
		ids, err := runner.Invoke(ctx, document.Source{URI: path})
		if err != nil {
			return fmt.Errorf("invoke index graph failed: %w", err)
		}

		fmt.Printf("[done] indexing file: %s, len of parts: %d\n", path, len(ids))

		return nil
	})

	return err
}

type RedisVectorStoreConfig struct {
	RedisKeyPrefix string
	IndexName      string
	Embedding      embedding.Embedder
	Dimension      int
	MinScore       float64
	RedisAddr      string
}

// initVectorIndex 初始化Redis中的向量索引
// 参数:
//
//	ctx: 上下文，用于取消操作和传递请求范围的值
//	config: RedisVectorStoreConfig的指针，包含连接和索引配置
//
// 返回值:
//
//	可能发生的错误
func initVectorIndex(ctx context.Context, config *RedisVectorStoreConfig) (err error) {
	// 检查embedding配置是否为空
	if config.Embedding == nil {
		return fmt.Errorf("embedding cannot be nil")
	}
	// 检查维度配置是否为正数
	if config.Dimension <= 0 {
		return fmt.Errorf("dimension must be positive")
	}

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})

	// 确保在错误时关闭连接
	defer func() {
		if err != nil {
			client.Close()
		}
	}()

	// 检查与Redis的连接
	if err = client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// 构建索引名称
	indexName := fmt.Sprintf("%s%s", config.RedisKeyPrefix, config.IndexName)

	// 检查是否存在索引
	exists, err := client.Do(ctx, "FT.INFO", indexName).Result()
	if err != nil {
		if !strings.Contains(err.Error(), "Unknown index name") {
			return fmt.Errorf("failed to check if index exists: %w", err)
		}
		err = nil
	} else if exists != nil {
		return nil
	}

	// 创建新索引
	createIndexArgs := []interface{}{
		"FT.CREATE", indexName,
		"ON", "HASH",
		"PREFIX", "1", config.RedisKeyPrefix,
		"SCHEMA",
		"content", "TEXT",
		"metadata", "TEXT",
		"vector", "VECTOR", "FLAT",
		"6",
		"TYPE", "FLOAT32",
		"DIM", config.Dimension,
		"DISTANCE_METRIC", "COSINE",
	}

	// 执行索引创建操作
	if err = client.Do(ctx, createIndexArgs...).Err(); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// 验证索引是否创建成功
	if _, err = client.Do(ctx, "FT.INFO", indexName).Result(); err != nil {
		return fmt.Errorf("failed to verify index creation: %w", err)
	}

	// 索引成功创建
	return nil
}
