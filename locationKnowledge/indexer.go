package locationKnowledge

import (
	"context"
	"encoding/json"
	"fmt"

	redispkg "my_assistant/pkg/redis"

	"log"
	"os"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	redisCli "github.com/redis/go-redis/v9"

	"github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino/components/indexer"
)

// 包初始化时自动调用，初始化 Redis 索引
func init() {
	err := redispkg.Init()
	if err != nil {
		log.Fatalf("failed to init redis index: %v", err)
	}
}

// newIndexer 是 locationKnowledge 图中节点 'Indexer1' 的组件初始化函数
// 该函数用于创建并配置一个基于 Redis 的向量索引器
func newIndexer(ctx context.Context) (idr indexer.Indexer, err error) {
	// 从环境变量获取 Redis 地址
	redisAddr := os.Getenv("REDIS_ADDR")
	// 创建 Redis 客户端
	redisClient := redisCli.NewClient(&redisCli.Options{
		Addr:     redisAddr,
		Protocol: 2,
	})

	// 配置索引器参数
	config := &redis.IndexerConfig{
		Client:    redisClient,          // Redis 客户端
		KeyPrefix: redispkg.PrefixRedis, // 键前缀
		BatchSize: 1,                    // 批量写入大小
		// DocumentToHashes 用于将 schema.Document 转换为 Redis Hashes
		DocumentToHashes: func(ctx context.Context, doc *schema.Document) (*redis.Hashes, error) {
			if doc.ID == "" {
				doc.ID = uuid.New().String() // 若无ID则生成UUID
			}
			key := doc.ID

			metadataBytes, err := json.Marshal(doc.MetaData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal metadata: %w", err)
			}

			return &redis.Hashes{
				Key: key,
				Field2Value: map[string]redis.FieldValue{
					redispkg.ContentField:  {Value: doc.Content, EmbedKey: redispkg.VectorField}, // 内容字段及向量字段
					redispkg.MetadataField: {Value: metadataBytes},                               // 元数据字段
				},
			}, nil
		},
	}

	// 初始化嵌入模型实例
	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	// 创建 Redis 索引器
	idr, err = redis.NewIndexer(ctx, config)
	if err != nil {
		return nil, err
	}
	return idr, nil
}
