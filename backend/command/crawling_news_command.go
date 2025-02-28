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
	"github.com/mjiee/world-news/backend/pkg/urlx"
	"github.com/mjiee/world-news/backend/service"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
)

// CrawlingNewsCommand is a command for crawling news.
type CrawlingNewsCommand struct {
	ctx context.Context

	crawlingSvc     service.CrawlingService
	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

func NewCrawlingNewsCommand(ctx context.Context, crawlingSvc service.CrawlingService, newsSvc service.NewsService,
	systemConfigSvc service.SystemConfigService) *CrawlingNewsCommand {
	return &CrawlingNewsCommand{
		ctx:             ctx,
		crawlingSvc:     crawlingSvc,
		newsSvc:         newsSvc,
		systemConfigSvc: systemConfigSvc,
	}
}

func (c *CrawlingNewsCommand) Execute(ctx context.Context) error {
	// check crawling record
	if err := c.crawlingNewsAllowed(ctx); err != nil {
		return err
	}

	// get news website
	newsWebsites, err := c.getNewsWebsites(ctx)
	if err != nil {
		return err
	}

	// get news keywords
	newsTopics, err := c.getNewsTopics(ctx)
	if err != nil {
		return err
	}

	// create crawling record
	record := entity.NewCrawlingRecord(valueobject.CrawlingNews,
		valueobject.NewCrawlingRecordConfig(newsWebsites, newsTopics))

	if err := c.crawlingSvc.CreateCrawlingRecord(ctx, record); err != nil {
		return err
	}

	// crawling news website
	go c.crawlingHandle(record)

	return nil
}

// crawlingNewsAllowed check crawling news allowed
func (c *CrawlingNewsCommand) crawlingNewsAllowed(ctx context.Context) error {
	// check crawling record
	hasProcessingTask, err := c.crawlingSvc.HasProcessingTasks(ctx)
	if err != nil {
		return err
	}

	if hasProcessingTask {
		return errorx.HasProcessingTasks
	}

	return nil
}

// getNewsWebsites get news websites
func (c *CrawlingNewsCommand) getNewsWebsites(ctx context.Context) ([]*valueobject.NewsWebsite, error) {
	websiteConfig, err := c.systemConfigSvc.GetSystemConfig(ctx, valueobject.NewsWebsiteKey.String())
	if err != nil {
		return nil, err
	}

	if websiteConfig.Id == 0 {
		return nil, errorx.NewsWebsiteConfigNotFound
	}

	newsWebsites, ok := websiteConfig.Value.([]*valueobject.NewsWebsite)
	if !ok {
		return nil, errorx.InternalError.SetErr(errors.New("invalid news websites config"))
	}

	return newsWebsites, nil
}

// getNewsTopics get news topics
func (c *CrawlingNewsCommand) getNewsTopics(ctx context.Context) ([]string, error) {
	topicConfig, err := c.systemConfigSvc.GetSystemConfig(ctx, valueobject.NewsTopicKey.String())
	if err != nil {
		return nil, err
	}

	if topicConfig.Id == 0 {
		return nil, errorx.NewsTopicConfigNotFound
	}

	newsTopics, ok := topicConfig.Value.([]string)
	if !ok {
		return nil, errorx.InternalError.SetErr(errors.New("invalid news topic config"))
	}

	return newsTopics, nil
}

// crawlingHandle crawling news
func (c *CrawlingNewsCommand) crawlingHandle(record *entity.CrawlingRecord) {
	for _, website := range record.Config.Sources {
		select {
		case <-c.ctx.Done():
			record.CrawlingPaused()

			return
		default:
			// crawling news
			newsQuantity, err := c.crawlingNews(website, record)
			if err != nil {
				logx.WithContext(c.ctx).Error("crawlingNews", err)

				return
			}

			// update crawling record quantity
			record, err = c.crawlingSvc.GetCrawlingRecord(c.ctx, record.Id)
			if err != nil {
				logx.WithContext(c.ctx).Error("GetCrawlingRecord", err)

				return
			}

			record.Quantity += newsQuantity

			if err := c.crawlingSvc.UpdateCrawlingRecordQuantity(c.ctx, record.Id, record.Quantity); err != nil {
				logx.WithContext(c.ctx).Error("UpdateCrawlingRecord", err)

				return
			}

			// check crawling record status, if not processing, interrupt handling
			if !record.Status.IsProcessing() {
				return
			}
		}
	}

	// update crawling record status
	record.CrawlingCompleted()

	if err := c.crawlingSvc.UpdateCrawlingRecord(c.ctx, record); err != nil {
		logx.WithContext(c.ctx).Error("UpdateCrawlingRecord", err)
	}
}

// crawlingNews crawling news
func (c *CrawlingNewsCommand) crawlingNews(website *valueobject.NewsWebsite, record *entity.CrawlingRecord) (int64, error) {
	// crawling news topic page
	topicPageUrls, err := c.crawlingNewsTopicPage(website, record.Config.Topics)
	if err != nil {
		logx.WithContext(c.ctx).Error("crawlingNewsTopicPage", err)

		return 0, nil
	}

	// crawling newsData
	newsData := c.crawlingNewsInTopicPage(topicPageUrls, website.Selector)

	// crawling news detail
	newsDetails := make([]*entity.NewsDetail, 0, len(newsData))

	for _, detail := range newsData {
		detail.RecordId = record.Id
		detail.Source = website.GetHost()

		if err := c.crawlingNewsDetail(detail, website.Selector); err != nil {
			logx.WithContext(c.ctx).Error("crawlingNewsDetail", err)

			continue
		}

		newsDetails = append(newsDetails, detail)
	}

	// remove duplicate
	removeDuplicateImage(newsDetails)

	// save news
	if err := c.newsSvc.CreateNews(c.ctx, newsDetails...); err != nil {
		return 0, err
	}

	return int64(len(newsDetails)), nil
}

// crawlingNewsTopicPage crawling news topic page
func (c *CrawlingNewsCommand) crawlingNewsTopicPage(website *valueobject.NewsWebsite, topics []string) (
	map[string][]string, error) {
	var (
		collector     = c.crawlingSvc.GetCollector()
		topicPageData = map[string][]string{}
		selector      = valueobject.Html_a
	)

	if website.Selector != nil && website.Selector.Topic != "" {
		selector = website.Selector.Topic
	}

	collector.OnHTML(selector, func(h *colly.HTMLElement) {
		for _, topic := range topics {
			if !strings.EqualFold(h.Text, topic) {
				continue
			}

			link := urlx.UrlPrefixHandle(h.Attr(valueobject.Attr_href), h.Request.URL)

			if len(link) == 0 {
				continue
			}

			topicPageData[topic] = append(topicPageData[topic], link)
		}
	})

	if err := collector.Visit(website.Url); err != nil {
		if pkgCollector.IgnorableError(err) {
			return topicPageData, nil
		}

		return nil, errors.WithStack(err)
	}

	return topicPageData, nil
}

// crawlingNewsInTopicPage crawling news in topic page
func (c *CrawlingNewsCommand) crawlingNewsInTopicPage(topicPages map[string][]string,
	selector *valueobject.Selector) []*entity.NewsDetail {
	var (
		isVisited = map[string]struct{}{}
		news      = []*entity.NewsDetail{}
	)

	for topic, urls := range topicPages {
		for _, pageUrl := range urls {
			if _, ok := isVisited[pageUrl]; ok {
				continue
			}

			isVisited[pageUrl] = struct{}{}

			data, err := c.crawlingNewsLink(pageUrl, topic, selector)
			if err != nil {
				logx.WithContext(c.ctx).Error("crawlingNewsLink", err)

				continue
			}

			news = append(news, data...)
		}
	}

	// remove duplicate
	news = slices.CompactFunc(news, func(a, b *entity.NewsDetail) bool {
		return a.Link == b.Link
	})

	return news
}

// crawlingNewsLink crawling news link
func (c *CrawlingNewsCommand) crawlingNewsLink(pageUrl, topic string, selector *valueobject.Selector) (
	[]*entity.NewsDetail, error) {
	var (
		collector    = c.crawlingSvc.GetCollector()
		news         = []*entity.NewsDetail{}
		linkSelector = valueobject.Html_a
	)

	if selector != nil && selector.Link != "" {
		linkSelector = selector.Link
	}

	collector.OnHTML(linkSelector, func(h *colly.HTMLElement) {
		headers := strings.Fields(h.Text)

		if len(headers) < 5 { // news title length must be greater than 5
			return
		}

		link := urlx.UrlPrefixHandle(h.Attr(valueobject.Attr_href), h.Request.URL)

		if len(link) == 0 {
			return
		}

		news = append(news, &entity.NewsDetail{
			Link:  link,
			Topic: topic,
		})
	})

	if err := collector.Visit(pageUrl); err != nil {
		if pkgCollector.IgnorableError(err) {
			return news, nil
		}

		return nil, errors.WithStack(err)
	}

	return news, nil
}

// crawlingNewsDetail crawling news detail
func (c *CrawlingNewsCommand) crawlingNewsDetail(news *entity.NewsDetail, selector *valueobject.Selector) error {
	var (
		collector = c.crawlingSvc.GetCollector()
	)

	// publish time
	collector.OnHTML(news.ExtractPublishTime(selector))

	// title
	collector.OnHTML(news.ExtractTitle(selector))

	// content
	collector.OnHTML(news.ExtractContent(selector))

	// images
	collector.OnHTML(news.ExtractImage(selector))

	if err := collector.Visit(news.Link); err != nil && !pkgCollector.IgnorableError(err) {
		return errors.WithStack(err)
	}

	// validate news
	return news.Validate()
}

// removeDuplicateImage removes duplicate elements.
func removeDuplicateImage(data []*entity.NewsDetail) {
	// count image
	imageCount := map[string]int{}

	for _, item := range data {
		for _, image := range item.Images {
			imageCount[image]++
		}
	}

	// remove duplicate image
	for _, item := range data {
		newImages := make([]string, 0, len(item.Images))

		for _, image := range item.Images {
			if imageCount[image] > 1 {
				continue
			}

			newImages = append(newImages, image)
		}

		item.Images = newImages
	}
}
