package command

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/urlx"
	"github.com/mjiee/world-news/backend/service"
)

// CrawlingNewsWebsiteCommand crawlling news website command
type CrawlingNewsWebsiteCommand struct {
	ctx context.Context

	crawlingSvc     service.CrawlingService
	systemConfigSvc service.SystemConfigService
}

func NewCrawlingNewsWebsiteCommand(ctx context.Context, crawlingSvc service.CrawlingService,
	systemConfigSvc service.SystemConfigService) *CrawlingNewsWebsiteCommand {
	return &CrawlingNewsWebsiteCommand{
		ctx:             ctx,
		crawlingSvc:     crawlingSvc,
		systemConfigSvc: systemConfigSvc,
	}
}

func (c *CrawlingNewsWebsiteCommand) Execute(ctx context.Context) error {
	// check crawling record
	hasProcessingTask, err := c.crawlingSvc.HasProcessingTasks(ctx)
	if err != nil {
		return err
	}

	if hasProcessingTask {
		return errorx.HasProcessingTasks
	}

	// get news website collection
	config, err := c.systemConfigSvc.GetSystemConfig(ctx, valueobject.NewsWebsiteCollectionKey.String())
	if err != nil {
		return err
	}

	if config.Id == 0 {
		return errorx.NewsWebsiteConfigNotFound
	}

	var newsWebsiteCollections []*valueobject.NewsWebsite
	if err := config.UnmarshalValue(&newsWebsiteCollections); err != nil {
		return errorx.InternalError.SetErr(errors.New("invalid news websites config"))
	}

	// create crawling record
	record := entity.NewCrawlingRecord(valueobject.CrawlingWebsite,
		valueobject.NewCrawlingRecordConfig(newsWebsiteCollections, nil))

	if err := c.crawlingSvc.CreateCrawlingRecord(ctx, record); err != nil {
		return err
	}

	// crawling news website
	go c.crawlingHandle(record)

	return nil
}

// crawlingHandle crawling news website
func (c *CrawlingNewsWebsiteCommand) crawlingHandle(record *entity.CrawlingRecord) {
	var (
		invalidNewsWebsites []string
		startTime           = time.Now()
	)

	// get news websites
	newsWebsites, err := c.systemConfigSvc.GetNewsWebsites(c.ctx)
	if err != nil {
		logx.WithContext(c.ctx).Error("GetNewsWebsite", err)
	}

	// get invalid news websites
	invalidNewsWebsitesConfig, err := c.systemConfigSvc.GetSystemConfig(c.ctx, valueobject.InvalidNewsWebsiteKey.String())
	if err != nil {
		logx.WithContext(c.ctx).Error("GetSystemConfig.InvalidNewsWebsiteKey", err)
	}

	if invalidNewsWebsitesConfig.Id != 0 {
		if err := invalidNewsWebsitesConfig.UnmarshalValue(&invalidNewsWebsites); err != nil {
			logx.WithContext(c.ctx).Error("UnmarshalValue.InvalidNewsWebsiteKey", err)
		}
	}

	// crawling news website
	for _, item := range record.Config.Sources {
		select {
		case <-c.ctx.Done():
			record.CrawlingPaused()

			logx.WithContext(c.ctx).Info("crawlingHandle", "crawling news website paused")

			_ = c.saveCrawlingResults(record, newsWebsites, invalidNewsWebsites)

			return
		default:
			// crawling news website
			websites, err := c.crawlingNewsWebsite(item.Url, item.Selector, make(map[string]bool))
			if err != nil {
				logx.WithContext(c.ctx).Error("crawlingNewsWebsite: "+item.Url, err)

				continue
			}

			websites = gokit.SliceDistinct(websites, func(v *valueobject.NewsWebsite) string { return v.GetHost() })

			count := 0

			for idx, website := range websites {
				count++

				// check invalid news site
				if c.isInvalidateNewsSite(item, newsWebsites, invalidNewsWebsites)(website) {
					invalidNewsWebsites = append(invalidNewsWebsites, website.GetHost())
				} else {
					newsWebsites = append(newsWebsites, website)
				}

				if count < 20 && idx < len(websites)-1 {
					continue
				}

				// remove duplicate news website
				newsWebsites = gokit.SliceDistinct(newsWebsites, func(v *valueobject.NewsWebsite) string { return v.GetHost() })

				// update record
				record, err = c.crawlingSvc.GetCrawlingRecord(c.ctx, record.Id)
				if err != nil {
					logx.WithContext(c.ctx).Error("crawlingHandle.GetCrawlingRecord", err)

					return
				}

				record.Quantity = int64(len(newsWebsites))

				if time.Since(startTime) > maxWorkTime || idx == len(websites)-1 {
					record.CrawlingCompleted()
				}

				// save crawling results
				if err = c.saveCrawlingResults(record, newsWebsites, invalidNewsWebsites); err != nil {
					logx.WithContext(c.ctx).Error("crawlingHandle.saveCrawlingResults", err)

					return
				}

				if !record.Status.IsProcessing() {
					return
				}

				count = 0
			}
		}
	}

	logx.WithContext(c.ctx).Info("crawlingHandle", fmt.Sprintf("crawling news website completed, quantity: %d",
		record.Quantity))
}

// saveCrawlingResults save crawling results
func (c *CrawlingNewsWebsiteCommand) saveCrawlingResults(record *entity.CrawlingRecord,
	newsWebsites []*valueobject.NewsWebsite, invalidNewsWebsites []string) error {
	// save news website
	if err := c.systemConfigSvc.SaveNewsWebsites(c.ctx, newsWebsites); err != nil {
		return err
	}

	// invalid news website
	invalidNewsWebsitesConfig, err := entity.NewSystemConfig(valueobject.InvalidNewsWebsiteKey.String(), invalidNewsWebsites)
	if err != nil {
		return err
	}

	if err := c.systemConfigSvc.SaveSystemConfig(c.ctx, invalidNewsWebsitesConfig); err != nil {
		return err
	}

	// update crawling record
	if err := c.crawlingSvc.UpdateCrawlingRecord(c.ctx, record); err != nil {
		return err
	}

	return nil
}

// crawlingNewsWebsite crawling news website
func (c *CrawlingNewsWebsiteCommand) crawlingNewsWebsite(source string, selector *valueobject.Selector,
	visited map[string]bool) ([]*valueobject.NewsWebsite, error) {
	var newsWebsites []*valueobject.NewsWebsite

	if selector == nil || visited[source] {
		return newsWebsites, nil
	}

	visited[source] = true

	// crawling website
	collector := c.crawlingSvc.GetCollector()

	collector.OnHTML(selector.Website, func(h *colly.HTMLElement) {
		link := h.Attr(valueobject.Attr_href)

		if urlx.IsValidURL(link) {
			newsWebsites = append(newsWebsites, &valueobject.NewsWebsite{Url: urlx.NormalizeURL(source, link)})
		}
	})

	if err := collector.Visit(source); err != nil {
		return nil, errors.WithStack(err)
	}

	if selector.Child == nil || len(newsWebsites) == 0 {
		return newsWebsites, nil
	}

	// crawling sub website
	newsWebsitesData := []*valueobject.NewsWebsite{}

	for _, link := range newsWebsites {
		if urlx.ExtractHostFromURL(source) != link.GetHost() {
			continue
		}

		data, err := c.crawlingNewsWebsite(link.Url, selector.Child, visited)
		if err != nil {
			logx.WithContext(c.ctx).Error("crawlingChildNewsWebsite:"+link.Url, err)

			continue
		}

		newsWebsitesData = append(newsWebsitesData, data...)
	}

	return append(newsWebsitesData, newsWebsites...), nil
}

// isInvalidateNewsSite check news site is invalidate
func (c *CrawlingNewsWebsiteCommand) isInvalidateNewsSite(source *valueobject.NewsWebsite, exists []*valueobject.NewsWebsite,
	invalidData []string) func(v *valueobject.NewsWebsite) bool {
	return func(v *valueobject.NewsWebsite) bool {
		if source.GetHost() == v.GetHost() {
			return true
		}

		if slices.Contains(invalidData, v.GetHost()) {
			return true
		}

		if slices.ContainsFunc(exists, func(e *valueobject.NewsWebsite) bool { return e.GetHost() == v.GetHost() }) {
			return true
		}

		newsCmd := NewCrawlingNewsCommand(c.ctx, time.Now().Add(-valueobject.MaxValidityPeriod).Format(time.DateTime),
			nil, nil, c.crawlingSvc, nil, c.systemConfigSvc)

		news, err := newsCmd.extractNewsList(0, valueobject.NewNewsTopicLink("", v.Url))
		if err != nil {
			logx.WithContext(c.ctx).Error(fmt.Sprintf("isInvalidateNewsSite.extractNewsList:%s", v.Url), err)

			return true
		}

		return len(news) == 0
	}
}
