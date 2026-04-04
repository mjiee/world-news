package openai

import (
	"context"

	eino_openai "github.com/cloudwego/eino-ext/components/model/openai"
)

// NewChatModel creates a new chat model
func NewChatModel(ctx context.Context, config *Config) *eino_openai.ChatModel {
	model, _ := eino_openai.NewChatModel(ctx, &eino_openai.ChatModelConfig{
		APIKey:  config.ApiKey,
		BaseURL: config.ApiUrl,
		Model:   config.Model,
	})

	return model
}
