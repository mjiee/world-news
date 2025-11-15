package command

import (
	"context"
	"strings"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/openai"
	"github.com/mjiee/world-news/backend/service"
)

// CritiqueNewsCommand represents a command for news critique.
type CritiqueNewsCommand struct {
	contents        []string
	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

func NewCritiqueNewsCommand(contents []string, newsSvc service.NewsService,
	systemConfigSvc service.SystemConfigService) *CritiqueNewsCommand {
	return &CritiqueNewsCommand{
		contents:        contents,
		newsSvc:         newsSvc,
		systemConfigSvc: systemConfigSvc,
	}
}

// newsCritiquePromptConfig represents the configuration for news critique prompt.
type newsCritiquePromptConfig struct {
	SystemPrompt string `json:"systemPrompt"`
}

func (c CritiqueNewsCommand) Execute(ctx context.Context) ([]string, error) {
	if len(c.contents) == 0 {
		return nil, errorx.ParamsError
	}

	// get openai config
	openaiConfig, err := c.systemConfigSvc.GetSystemConfig(ctx, valueobject.TextAIKey.String())
	if err != nil {
		return nil, err
	}

	textAiConfig, err := entity.UnmarshalValue[openai.Config](openaiConfig, errorx.OpenaiConfigNotFound)
	if err != nil {
		return nil, err
	}

	// get news ceritique prompt
	critiquePrompt, err := c.systemConfigSvc.GetSystemConfig(ctx, valueobject.NewsCritiquePromptKey.String())
	if err != nil {
		return nil, err
	}

	critiqueConfig, err := entity.UnmarshalValue[newsCritiquePromptConfig](critiquePrompt, errorx.CritiquePromptNotFound)
	if err != nil {
		return nil, err
	}

	// news critique
	data, err := openai.NewOpenaiClient(textAiConfig).SetSystemPrompt(critiqueConfig.SystemPrompt).
		SetUserPrompt(c.contents...).ChatCompletion(ctx)
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
