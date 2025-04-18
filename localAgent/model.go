package localAgent

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
)

func newChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	config := &ark.ChatModelConfig{
		Model:  os.Getenv("ARK_CHAT_MODEL"),
		APIKey: os.Getenv("ARK_API_KEY"),
	}
	cm, err = ark.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
