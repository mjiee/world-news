package command

import (
	"context"
	"slices"
	"strings"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	pkgCollector "github.com/mjiee/world-news/backend/pkg/collector"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/service"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
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

	newsWebsiteCollections, ok := config.Value.([]*valueobject.NewsWebsite)
	if !ok {
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
		newsWebsites = make([]*valueobject.NewsWebsite, 0)
		err          error
	)

	for _, item := range record.Config.Sources {
		select {
		case <-c.ctx.Done():
			record.CrawlingPaused()

			return
		default:
			// check crawling record
			record, err = c.crawlingSvc.GetCrawlingRecord(c.ctx, record.Id)
			if err != nil {
				logx.WithContext(c.ctx).Error("GetCrawlingRecord", err)

				record.CrawlingFailed()

				return
			}

			if !record.Status.IsProcessing() {
				return
			}

			// crawling news website
			websites, err := c.crawlingNewsWebsite(item.Url, item.Selector)
			if err != nil {
				logx.WithContext(c.ctx).Error("crawlingNewsWebsite", err)

				continue
			}

			websites = slices.DeleteFunc(websites, func(v *valueobject.NewsWebsite) bool {
				return strings.HasPrefix(v.Url, item.Url)
			})

			newsWebsites = append(newsWebsites, websites...)
		}
	}

	// remove duplicate
	newsWebsites = slices.CompactFunc(newsWebsites, func(a, b *valueobject.NewsWebsite) bool {
		return a.GetHost() == b.GetHost()
	})

	// save news website
	if err := c.systemConfigSvc.SaveSystemConfig(c.ctx,
		entity.NewSystemConfig(valueobject.NewsWebsiteKey.String(), newsWebsites)); err != nil {

		logx.WithContext(c.ctx).Error("SaveSystemConfig", err)

		return
	}

	// update crawling record
	record.CrawlingCompleted()
	record.Quantity = int64(len(newsWebsites))

	if err := c.crawlingSvc.UpdateCrawlingRecord(c.ctx, record); err != nil {
		logx.WithContext(c.ctx).Error("SaveSystemConfig", err)
	}
}

// crawlingNewsWebsite crawling news website
func (c *CrawlingNewsWebsiteCommand) crawlingNewsWebsite(collectionUrl string,
	selector *valueobject.Selector) ([]*valueobject.NewsWebsite, error) {

	var newsWebsites []*valueobject.NewsWebsite

	if selector == nil {
		return newsWebsites, nil
	}

	// crawling website
	collector := c.crawlingSvc.GetCollector()

	collector.OnHTML(selector.Website, func(h *colly.HTMLElement) {
		link := h.Attr(valueobject.Attr_href)

		if len(link) > 0 {
			newsWebsites = append(newsWebsites, &valueobject.NewsWebsite{Url: link})
		}
	})

	if err := collector.Visit(collectionUrl); err != nil {
		if pkgCollector.IgnorableError(err) {
			return newsWebsites, nil
		}

		return nil, errors.WithStack(err)
	}

	if selector.Child == nil || len(newsWebsites) == 0 {
		return newsWebsites, nil
	}

	// crawling sub website
	newsWebsitesData := []*valueobject.NewsWebsite{}

	for _, link := range newsWebsites {
		data, err := c.crawlingNewsWebsite(link.Url, selector.Child)
		if err != nil {
			return nil, err
		}

		newsWebsitesData = append(newsWebsitesData, data...)
	}

	return append(newsWebsitesData, newsWebsites...), nil
}
