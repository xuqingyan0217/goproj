package redis

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
)

// Redis 相关常量定义
const (
	PrefixRedis = "excel:doc:"   // Redis 键前缀
	IndexName   = "vector_index" // Redis 索引名称

	ContentField  = "content"        // 内容字段名
	MetadataField = "metadata"       // 元数据字段名
	VectorField   = "content_vector" // 向量字段名
	DistanceField = "distance"       // 距离字段名
)

// 用于保证只初始化一次的 sync.Once
var initOnce sync.Once

// Init 初始化 Redis 索引（只执行一次）
func Init() error {
	var err error
	initOnce.Do(func() {
		err = InitRedisIndex(context.Background(), &Config{
			RedisAddr: os.Getenv("REDIS_ADDR"), // 从环境变量获取 Redis 地址
			Dimension: 2560,                    // 向量维度
		})
	})
	return err
}

// Config Redis 索引初始化配置
// RedisAddr: Redis 服务地址
// Dimension: 向量维度
type Config struct {
	RedisAddr string
	Dimension int
}

// InitRedisIndex 初始化 Redis 向量索引
// 若索引不存在则创建，存在则跳过
func InitRedisIndex(ctx context.Context, config *Config) (err error) {
	if config.Dimension <= 0 {
		return fmt.Errorf("dimension must be positive")
	}

	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Protocol: 2,
	})

	// 若出错则关闭客户端
	defer func() {
		if err != nil {
			_ = client.Close()
		}
	}()

	// 检查 Redis 连接
	if err = client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	indexName := fmt.Sprintf("%s%s", PrefixRedis, IndexName)

	// 检查索引是否存在
	exists, err := client.Do(ctx, "FT.INFO", indexName).Result()
	if err != nil {
		if !strings.Contains(err.Error(), "Unknown index name") {
			return fmt.Errorf("failed to check if index exists: %w", err)
		}
		err = nil // 未找到索引则继续
	} else if exists != nil {
		return nil // 索引已存在，无需创建
	}

	// 构造创建索引的参数
	createIndexArgs := []interface{}{
		"FT.CREATE", indexName,
		"ON", "HASH",
		"PREFIX", "1", PrefixRedis,
		"SCHEMA",
		ContentField, "TEXT",
		MetadataField, "TEXT",
		VectorField, "VECTOR", "FLAT",
		"6",
		"TYPE", "FLOAT32",
		"DIM", config.Dimension,
		"DISTANCE_METRIC", "COSINE",
	}

	// 创建新索引
	if err = client.Do(ctx, createIndexArgs...).Err(); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// 验证索引是否创建成功
	if _, err = client.Do(ctx, "FT.INFO", indexName).Result(); err != nil {
		return fmt.Errorf("failed to verify index creation: %w", err)
	}

	return nil
}
