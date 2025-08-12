package command

import (
	"context"

	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/translate"
	"github.com/mjiee/world-news/backend/service"

	"github.com/pkg/errors"
)

// TranslateNewsCommand represents a command for news translation.
type TranslateNewsCommand struct {
	contents []string
	toLang   string

	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

func NewTranslateNewsCommand(contents []string, toLang string, newsSvc service.NewsService,
	systemConfigSvc service.SystemConfigService) *TranslateNewsCommand {
	return &TranslateNewsCommand{
		contents:        contents,
		toLang:          toLang,
		newsSvc:         newsSvc,
		systemConfigSvc: systemConfigSvc}
}

func (c TranslateNewsCommand) Execute(ctx context.Context) ([]string, error) {
	if len(c.contents) == 0 {
		return nil, errorx.ParamsError
	}

	// get translater config
	translaterConfig, err := c.systemConfigSvc.GetSystemConfig(ctx, valueobject.TranslaterKey.String())
	if err != nil {
		return nil, err
	}

	if translaterConfig.Id == 0 {
		return nil, errorx.TranslaterConfigNotFound
	}

	var config translate.Config
	if err := translaterConfig.UnmarshalValue(&config); err != nil {
		return nil, errorx.InternalError.SetErr(errors.New("translater config type error"))
	}

	// translate
	translater, err := translate.NewTranslater(&config)
	if err != nil {
		return nil, err
	}

	return translater.Translate(ctx, c.toLang, c.contents...)
}
