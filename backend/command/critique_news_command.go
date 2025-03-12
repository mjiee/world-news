package command

import (
	"context"
	"errors"
	"strings"

	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/openai"
	"github.com/mjiee/world-news/backend/service"
)

// CritiqueNewsCommand represents a command for news critique.
type CritiqueNewsCommand struct {
	newsId          uint
	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

func NewCritiqueNewsCommand(newsId uint, newsSvc service.NewsService,
	systemConfigSvc service.SystemConfigService) *CritiqueNewsCommand {
	return &CritiqueNewsCommand{
		newsId:          newsId,
		newsSvc:         newsSvc,
		systemConfigSvc: systemConfigSvc,
	}
}

func (c CritiqueNewsCommand) Execute(ctx context.Context) ([]string, error) {
	// get news
	news, err := c.newsSvc.GetNewsDetail(ctx, c.newsId)
	if err != nil {
		return nil, err
	}

	// get openai config
	openaiConfig, err := c.systemConfigSvc.GetSystemConfig(ctx, valueobject.OpenAIKey.String())
	if err != nil {
		return nil, err
	}

	if openaiConfig.Id == 0 {
		return nil, errorx.OpenaiConfigNotFound
	}

	config, ok := openaiConfig.Value.(openai.Config)
	if !ok {
		return nil, errorx.InternalError.SetErr(errors.New("openai config type error"))
	}

	// news critique
	data, err := openai.NewOpenaiClient(&config).ChatCompletion(ctx, strings.Join(news.Contents, "\n"))
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for _, choice := range data.Choices {
		contents := strings.Split(choice.Message.Content, "\n")

		result = append(result, contents...)
	}

	return result, nil
}
