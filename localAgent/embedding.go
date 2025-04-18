package localAgent

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/components/embedding"
)

func newEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	config := &ark.EmbeddingConfig{
		Model:  os.Getenv("ARK_EMBEDDING_MODEL"),
		APIKey: os.Getenv("ARK_API_KEY"),
	}
	eb, err = ark.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	return eb, nil
}
