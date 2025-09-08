package command

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/textx"
	"github.com/mjiee/world-news/backend/pkg/urlx"
	"github.com/mjiee/world-news/backend/service"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mjiee/gokit/slicex"
	"github.com/pkg/errors"
)

// maxWorkTime is the maximum time to work.
const maxWorkTime = 8 * time.Hour

// CrawlingNewsCommand is a command for crawling news.
type CrawlingNewsCommand struct {
	ctx       context.Context
	startTime time.Time
	sources   []string
	topics    []string

	crawlingSvc     service.CrawlingService
	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

func NewCrawlingNewsCommand(ctx context.Context, startTime string, sources []string, topics []string,
	crawlingSvc service.CrawlingService, newsSvc service.NewsService, systemConfigSvc service.SystemConfigService,
) *CrawlingNewsCommand {
	cmd := &CrawlingNewsCommand{
		ctx:             ctx,
		sources:         sources,
		topics:          topics,
		crawlingSvc:     crawlingSvc,
		newsSvc:         newsSvc,
		systemConfigSvc: systemConfigSvc,
	}

	if startTime != "" {
		cmd.startTime, _ = time.Parse(time.DateOnly, startTime)
	} else {
		cmd.startTime = time.Now().AddDate(0, 0, -1)
	}

	return cmd
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
	if err := c.getNewsTopics(ctx); err != nil {
		return err
	}

	// create crawling record
	record := entity.NewCrawlingRecord(valueobject.CrawlingNews,
		valueobject.NewCrawlingRecordConfig(newsWebsites, c.topics))

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
	newsWebsites, err := c.systemConfigSvc.GetNewsWebsites(ctx)
	if err != nil {
		return nil, err
	}

	if len(c.sources) > 0 {
		newsWebsites = slicex.Filter(newsWebsites, func(nw *valueobject.NewsWebsite) bool {
			return slices.Contains(c.sources, urlx.ExtractSecondLevelDomain(nw.Url))
		})
	}

	if len(newsWebsites) == 0 {
		return nil, errorx.NewsWebsiteConfigNotFound
	}

	return newsWebsites, nil
}

// getNewsTopics get news topics
func (c *CrawlingNewsCommand) getNewsTopics(ctx context.Context) error {
	if len(c.topics) > 0 {
		return nil
	}

	topicConfig, err := c.systemConfigSvc.GetSystemConfig(ctx, valueobject.NewsTopicKey.String())
	if err != nil {
		return err
	}

	if topicConfig.Id == 0 {
		return nil
	}

	var newsTopics []string

	if err := topicConfig.UnmarshalValue(&newsTopics); err != nil {
		return errorx.InternalError.SetErr(errors.New("invalid news topic config"))
	}

	c.topics = newsTopics

	return nil
}

// crawlingHandle crawling news
func (c *CrawlingNewsCommand) crawlingHandle(record *entity.CrawlingRecord) {
	startTime := time.Now()

	for idx, website := range record.Config.Sources {
		select {
		case <-c.ctx.Done():
			record.CrawlingPaused()

			_ = c.crawlingSvc.UpdateCrawlingRecord(c.ctx, record)

			logx.WithContext(c.ctx).Info("crawlingHandle", "crawling news website paused")

			return
		default:
			// crawling news
			newsQuantity, err := c.crawlingNews(website, record)
			if err != nil {
				logx.WithContext(c.ctx).Error("crawlingHandle.crawlingNews:"+website.Url, err)

				continue
			}

			// update crawling record
			record, err = c.crawlingSvc.GetCrawlingRecord(c.ctx, record.Id)
			if err != nil {
				logx.WithContext(c.ctx).Error("GetCrawlingRecord:", err)

				return
			}

			record.Quantity += newsQuantity

			if time.Since(startTime) > maxWorkTime || idx == len(record.Config.Sources) {
				record.CrawlingCompleted()
			}

			if err := c.crawlingSvc.UpdateCrawlingRecord(c.ctx, record); err != nil {
				logx.WithContext(c.ctx).Error("UpdateCrawlingRecord", err)

				return
			}

			// check crawling record status, if not processing, interrupt handling
			if !record.Status.IsProcessing() {
				return
			}
		}
	}

	logx.WithContext(c.ctx).Info("crawlingHandle", fmt.Sprintf("crawling news website completed, quantity: %d",
		record.Quantity))
}

// crawlingNews crawling news
func (c *CrawlingNewsCommand) crawlingNews(website *valueobject.NewsWebsite, record *entity.CrawlingRecord,
) (int64, error) {
	// crawling news topic page
	topicLinks, err := c.extractNewsTopicLinks(website, record)
	if err != nil {
		return 0, err
	}

	// crawling newsData
	newsData := c.crawlingNewsInTopicPage(record, topicLinks)

	// save news
	if err := c.newsSvc.CreateNews(c.ctx, newsData...); err != nil {
		return 0, err
	}

	return int64(len(newsData)), nil
}

// extractNewsTopicLinks extract news topic links
func (c *CrawlingNewsCommand) extractNewsTopicLinks(website *valueobject.NewsWebsite, record *entity.CrawlingRecord,
) ([]*valueobject.NewsTopicLink, error) {
	if len(record.Config.Topics) == 0 {
		return []*valueobject.NewsTopicLink{valueobject.NewNewsTopicLink("", website.Url)}, nil
	}

	var (
		collector = c.crawlingSvc.GetCollector()
		result    []*valueobject.NewsTopicLink
	)

	for _, selector := range valueobject.NewsTopicLinkSelectors {
		collector.OnHTML(selector, func(e *colly.HTMLElement) {
			href := e.Attr(valueobject.Attr_href)
			if !urlx.IsValidURL(href) {
				return
			}

			linkText := textx.CleanText(e.Text)
			if linkText == "" {
				linkText = textx.CleanText(e.Attr(valueobject.Attr_title))
			}

			topic, matches := textx.MatchesKeyword(linkText, record.Config.Topics)
			if !matches {
				return
			}

			data := valueobject.NewNewsTopicLink(topic, urlx.NormalizeURL(website.Url, href))

			if slices.ContainsFunc(result, data.Compare) {
				return
			}

			result = append(result, data)
		})
	}

	err := collector.Visit(website.Url)

	return slicex.Distinct(result, func(i *valueobject.NewsTopicLink) string { return i.URL }), errors.WithStack(err)
}

// crawlingNewsInTopicPage crawling news in topic page
func (c *CrawlingNewsCommand) crawlingNewsInTopicPage(record *entity.CrawlingRecord, topicLinks []*valueobject.NewsTopicLink,
) []*entity.NewsDetail {
	result := []*entity.NewsDetail{}

	for _, link := range topicLinks {
		newsList, err := c.extractNewsList(record.Id, link)
		if err != nil {
			logx.WithContext(c.ctx).Error("extractNewsList", err)

			continue
		}

		result = append(result, newsList...)
	}

	return removeDuplicateNews(result)
}

// removeDuplicateNews removes duplicate elements.
func removeDuplicateNews(data []*entity.NewsDetail) []*entity.NewsDetail {
	var (
		images  = make([]string, 0)
		newsMap = slicex.GroupBy(data, func(item *entity.NewsDetail) string { return item.Link })
		result  = make([]*entity.NewsDetail, 0, len(data))
	)

	for _, newsData := range newsMap {
		news := slices.MaxFunc(newsData, func(a, b *entity.NewsDetail) int { return a.Compare(b) })

		news.Images = slicex.Filter(news.Images, func(v string) bool {
			v = urlx.RemoveQueryParams(v)

			if slices.Contains(images, v) {
				return false
			}

			images = append(images, v)

			return true
		})

		result = append(result, news)
	}

	return result
}

// extractNewsList extract news list
func (c *CrawlingNewsCommand) extractNewsList(recordId uint, link *valueobject.NewsTopicLink,
) (result []*entity.NewsDetail, err error) {
	collector := c.crawlingSvc.GetCollector()

	collector.OnHTML(valueobject.Html, func(e *colly.HTMLElement) {
		items := c.findNewsItems(e.DOM)

		if len(items) == 0 {
			items = c.findNewsLink(e.DOM)
		}

		newsList := slicex.Map(items, func(item *goquery.Selection) *entity.NewsDetail {
			detail := entity.NewNewsDetailFromTopicLink(recordId, link)

			detail.ExtractTitle(item)
			detail.ExtractSummary(item)
			detail.ExtractLink(link.URL, item)
			detail.ExtractImages(item)
			detail.ExtractPublishTime(c.startTime, item)

			return detail
		})

		result = append(result, newsList...)
	})

	err = collector.Visit(link.URL)

	return slicex.Filter(result, func(v *entity.NewsDetail) bool { return v != nil && v.IsValid(c.startTime) }),
		errors.WithStack(err)
}

// findNewsItems find news items
func (c *CrawlingNewsCommand) findNewsItems(doc *goquery.Selection) []*goquery.Selection {
	var items []*goquery.Selection

	for _, selector := range valueobject.ExcludeSelectors {
		doc.Find(selector).Remove()
	}

	for _, selector := range valueobject.NewsItemSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			if s.Find(valueobject.LinkSelector).Length() > 0 {
				items = append(items, s)
			}
		})
	}

	return items
}

// findNewsLink find news from link
func (c *CrawlingNewsCommand) findNewsLink(doc *goquery.Selection) []*goquery.Selection {
	var items []*goquery.Selection

	doc.Find(valueobject.LinkSelector).Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr(valueobject.Attr_href)
		if valueobject.IsNewsLink(href) {
			items = append(items, s)
		}
	})

	return items
}
