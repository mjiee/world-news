package command

import (
	"context"
	"slices"
	"strings"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	pkgCollector "github.com/mjiee/world-news/backend/pkg/collector"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/service"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
)

// CrawlingNewsWebsiteCommand crawlling news website command
type CrawlingNewsWebsiteCommand struct {
	newsCrawlingSvc   service.NewsCrawlingService
	systemSettingsSvc service.SystemSettingsService
}

func NewCrawlingNewsWebsiteCommand(newsCrawlingSvc service.NewsCrawlingService,
	systemSettingsSvc service.SystemSettingsService) *CrawlingNewsWebsiteCommand {
	return &CrawlingNewsWebsiteCommand{
		newsCrawlingSvc:   newsCrawlingSvc,
		systemSettingsSvc: systemSettingsSvc,
	}
}

func (c *CrawlingNewsWebsiteCommand) Execute(ctx context.Context) error {
	// get news website collection
	config, err := c.systemSettingsSvc.GetSystemConfig(ctx, valueobject.NewsWebsiteCollectionKey)
	if err != nil {
		return err
	}

	if config.Id == 0 {
		return errorx.InternalError
	}

	newsWebsiteCollections := config.Value.([]*valueobject.NewsWebsite)

	// create crawling record
	record := entity.NewCrawlingRecord(valueobject.CrawlingWebsite)

	if err := c.newsCrawlingSvc.CreateCrawlingRecord(ctx, record); err != nil {
		return err
	}

	// crawling news website
	go c.crawlingNewsWebsite(record, newsWebsiteCollections)

	return nil
}

// crawlingNewsWebsite crawling news website
func (c *CrawlingNewsWebsiteCommand) crawlingNewsWebsite(record *entity.CrawlingRecord,
	data []*valueobject.NewsWebsite) {
	newsWebsites := make([]string, 0)

	for _, item := range data {
		websites, err := c.crawlingHandle(item.Url, item.Selectors...)
		if err != nil {
			// TODO: logging

			continue
		}

		websites = slices.DeleteFunc(websites, func(v string) bool {
			return strings.HasPrefix(v, item.Url)
		})

		newsWebsites = append(newsWebsites, websites...)
	}

	// remove duplicate
	newsWebsites = slices.CompactFunc(newsWebsites, strings.EqualFold)

	data = make([]*valueobject.NewsWebsite, len(newsWebsites))

	for i, item := range newsWebsites {
		data[i] = &valueobject.NewsWebsite{Url: item}
	}

	// save news website
	if err := c.systemSettingsSvc.SaveSystemConfig(context.Background(),
		entity.NewSystemConfig(valueobject.NewsWebsiteKey, data)); err != nil {
		// TODO: logging

		return
	}

	// update crawling record
	record.Status = valueobject.CompletedCrawlingRecord
	record.Quantity = int64(len(newsWebsites))

	if err := c.newsCrawlingSvc.UpdateCrawlingRecord(context.Background(), record); err != nil {
		// TODO: logging
		return
	}
}

// crawlingHandle crawling handle
func (c *CrawlingNewsWebsiteCommand) crawlingHandle(website string, selectors ...string) ([]string, error) {
	if len(selectors) == 0 {
		return nil, nil
	}

	// crawling website
	var (
		links     = make([]string, 0)
		collector = c.newsCrawlingSvc.GetCollector()
	)

	collector.OnHTML(selectors[0], func(h *colly.HTMLElement) {
		link := h.Attr(valueobject.Attr_href)

		if len(link) > 0 {
			links = append(links, link)
		}
	})

	if err := collector.Visit(website); err != nil {
		if pkgCollector.IgnorableError(err) {
			return links, nil
		}

		return nil, errors.WithStack(err)
	}

	if len(selectors) <= 1 || len(links) == 0 {
		return links, nil
	}

	// crawling sub website
	newUrls := []string{}

	for _, link := range links {
		urls, err := c.crawlingHandle(link, selectors[1:]...)
		if err != nil {
			return nil, err
		}

		newUrls = append(newUrls, urls...)
	}

	return newUrls, nil
}
