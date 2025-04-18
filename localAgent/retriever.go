package localAgent

import (
	"context"
	"fmt"
	redispkg "my_assistant/pkg/redis"

	"os"
	"strconv"

	"github.com/cloudwego/eino/schema"
	redisCli "github.com/redis/go-redis/v9"

	"github.com/cloudwego/eino-ext/components/retriever/redis"

	"github.com/cloudwego/eino/components/retriever"
)

// newRetriever 是 localAgent 图中节点 'RedisRetriever' 的组件初始化函数
// 该函数用于创建并配置一个基于 Redis 的向量检索器
func newRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {
	// 从环境变量获取 Redis 地址
	redisAddr := os.Getenv("REDIS_ADDR")
	// 创建 Redis 客户端
	redisClient := redisCli.NewClient(&redisCli.Options{
		Addr:     redisAddr,
		Protocol: 2,
	})
	// 配置检索器参数
	config := &redis.RetrieverConfig{
		Client:       redisClient,                                                                     // Redis 客户端
		Index:        fmt.Sprintf("%s%s", redispkg.PrefixRedis, redispkg.IndexName),                   // 索引名称
		Dialect:      2,                                                                               // Redis 语法版本
		ReturnFields: []string{redispkg.ContentField, redispkg.MetadataField, redispkg.DistanceField}, // 需要返回的字段
		TopK:         8,                                                                               // 检索返回的文档数量
		VectorField:  redispkg.VectorField,                                                            // 向量字段名
		// DocumentConverter 用于将 Redis 文档转换为 schema.Document
		DocumentConverter: func(ctx context.Context, doc redisCli.Document) (*schema.Document, error) {
			resp := &schema.Document{
				ID:       doc.ID,
				Content:  "",
				MetaData: map[string]any{},
			}
			for field, val := range doc.Fields {
				if field == redispkg.ContentField {
					resp.Content = val
				} else if field == redispkg.MetadataField {
					resp.MetaData[field] = val
				} else if field == redispkg.DistanceField {
					distance, err := strconv.ParseFloat(val, 64)
					if err != nil {
						continue
					}
					resp.WithScore(1 - distance) // 距离转为分数
				}
			}

			return resp, nil
		},
	}
	// 初始化嵌入模型实例
	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	// 创建 Redis 检索器
	rtr, err = redis.NewRetriever(ctx, config)
	if err != nil {
		return nil, err
	}
	return rtr, nil
}
