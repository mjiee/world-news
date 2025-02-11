package command

import (
	"context"

	"github.com/mjiee/world-news/backend/service"
)

// CrawlingNewsCommand is a command for crawling news.
type CrawlingNewsCommand struct {
	newsCrawlingSvc   service.NewsCrawlingService
	newsDetailSvc     service.NewsDetailService
	systemSettingsSvc service.SystemSettingsService
}

func NewCrawlingNewsCommand(newsCrawlingSvc service.NewsCrawlingService, newsDetailSvc service.NewsDetailService,
	systemSettingsSvc service.SystemSettingsService) *CrawlingNewsCommand {
	return &CrawlingNewsCommand{
		newsCrawlingSvc:   newsCrawlingSvc,
		newsDetailSvc:     newsDetailSvc,
		systemSettingsSvc: systemSettingsSvc,
	}
}

func (c *CrawlingNewsCommand) Execute(ctx context.Context) error {
	return nil
}
