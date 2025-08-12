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
	title    string
	contents []string

	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

func NewCritiqueNewsCommand(title string, contents []string, newsSvc service.NewsService,
	systemConfigSvc service.SystemConfigService) *CritiqueNewsCommand {
	return &CritiqueNewsCommand{
		title:           title,
		contents:        contents,
		newsSvc:         newsSvc,
		systemConfigSvc: systemConfigSvc,
	}
}

func (c CritiqueNewsCommand) Execute(ctx context.Context) ([]string, error) {
	if c.title == "" || len(c.contents) == 0 {
		return nil, errorx.ParamsError
	}

	// get openai config
	openaiConfig, err := c.systemConfigSvc.GetSystemConfig(ctx, valueobject.OpenAIKey.String())
	if err != nil {
		return nil, err
	}

	if openaiConfig.Id == 0 {
		return nil, errorx.OpenaiConfigNotFound
	}

	var config openai.Config
	if err := openaiConfig.UnmarshalValue(&config); err != nil {
		return nil, errorx.InternalError.SetErr(errors.New("openai config type error"))
	}

	// news critique
	data, err := openai.NewOpenaiClient(&config).ChatCompletion(ctx, c.title, strings.Join(c.contents, "\n"))
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
