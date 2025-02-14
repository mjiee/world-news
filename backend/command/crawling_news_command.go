package command

import (
	"context"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	pkgCollector "github.com/mjiee/world-news/backend/pkg/collector"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/service"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
)

// CrawlingNewsCommand is a command for crawling news.
type CrawlingNewsCommand struct {
	crawlingSvc     service.CrawlingService
	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

func NewCrawlingNewsCommand(crawlingSvc service.CrawlingService, newsSvc service.NewsService,
	systemConfigSvc service.SystemConfigService) *CrawlingNewsCommand {
	return &CrawlingNewsCommand{
		crawlingSvc:     crawlingSvc,
		newsSvc:         newsSvc,
		systemConfigSvc: systemConfigSvc,
	}
}

func (c *CrawlingNewsCommand) Execute(ctx context.Context) error {
	// get news website
	websiteConfig, err := c.systemConfigSvc.GetSystemConfig(ctx, valueobject.NewsWebsiteKey)
	if err != nil {
		return err
	}

	if websiteConfig.Id == 0 {
		return errorx.InternalError
	}

	newsWebsites := websiteConfig.Value.([]*valueobject.NewsWebsite)

	// get news keywords
	topicConfig, err := c.systemConfigSvc.GetSystemConfig(ctx, valueobject.NewsTopicKey)
	if err != nil {
		return err
	}

	if websiteConfig.Id == 0 {
		return errorx.InternalError
	}

	newsKeywords := topicConfig.Value.([]string)

	// create crawling record
	record := entity.NewCrawlingRecord(valueobject.CrawlingNews)

	if err := c.crawlingSvc.CreateCrawlingRecord(ctx, record); err != nil {
		return err
	}

	// crawling news website
	go c.crawlingNewsHandle(record, newsKeywords, newsWebsites)

	return nil
}

// crawlingNewsHandle crawling news
func (c *CrawlingNewsCommand) crawlingNewsHandle(record *entity.CrawlingRecord, newsTopics []string,
	newsWebsites []*valueobject.NewsWebsite) {
	for _, website := range newsWebsites {
		// crawling news topic page
		topicPageUrls, err := c.crawlingNewsTopicPage(website.Url, newsTopics)
		if err != nil {
			// TODO: logging

			continue
		}

		// crawling news link
		newsLinks := make([]string, 0)

		for _, topicPageUrl := range topicPageUrls {
			links, err := c.crawlingNewsLink(topicPageUrl)
			if err != nil {
				// TODO: logging
				continue
			}

			newsLinks = append(newsLinks, links...)
		}

		newsLinks = slices.CompactFunc(newsLinks, strings.EqualFold)

		// crawling news detail
		news := make([]*entity.NewsDetail, 0)

		for _, newsLink := range newsLinks {
			detail, err := c.crawlingNewsDetail(newsLink)
			if err != nil {
				// TODO: logging
				continue
			}

			if detail == nil {
				continue
			}

			news = append(news, detail)
		}

		// save news
		if err := c.newsSvc.CreateNews(context.Background(), news...); err != nil {
			// TODO: logging
			continue
		}

		record.Quantity += int64(len(news))
	}

	// update crawling record status
	record.Status = valueobject.CompletedCrawlingRecord

	if err := c.crawlingSvc.UpdateCrawlingRecord(context.Background(), record); err != nil {
		// TODO: logging
	}
}

// crawlingNewsTopicPage crawling news topic page
func (c *CrawlingNewsCommand) crawlingNewsTopicPage(websiteUrl string, topics []string) ([]string, error) {
	var (
		collector = c.crawlingSvc.GetCollector()
		pageUrls  = []string{}
	)

	collector.OnHTML(valueobject.Html_a, func(h *colly.HTMLElement) {
		if !isTopicLink(h.Text, topics) {
			return
		}

		link := urlPrefixHandle(h.Attr(valueobject.Attr_href), h.Request.URL)

		if len(link) == 0 {
			return
		}

		pageUrls = append(pageUrls, link)
	})

	if err := collector.Visit(websiteUrl); err != nil {
		if pkgCollector.IgnorableError(err) {
			return pageUrls, nil
		}

		return nil, errors.WithStack(err)
	}

	pageUrls = slices.CompactFunc(pageUrls, strings.EqualFold) // remove duplicate

	return pageUrls, nil
}

// crawlingNewsLink crawling news link
func (c *CrawlingNewsCommand) crawlingNewsLink(topicPageUrl string) ([]string, error) {
	var (
		collector = c.crawlingSvc.GetCollector()
		newsLinks = []string{}
	)

	collector.OnHTML(valueobject.Html_a, func(h *colly.HTMLElement) {
		headers := strings.Fields(h.Text)

		if len(headers) < 5 { // news title length must be greater than 5
			return
		}

		link := urlPrefixHandle(h.Attr(valueobject.Attr_href), h.Request.URL)

		if len(link) == 0 {
			return
		}

		newsLinks = append(newsLinks, link)
	})

	if err := collector.Visit(topicPageUrl); err != nil {
		if pkgCollector.IgnorableError(err) {
			return newsLinks, nil
		}

		return nil, errors.WithStack(err)
	}

	newsLinks = slices.CompactFunc(newsLinks, strings.EqualFold)

	if len(newsLinks) > 10 {
		return newsLinks[:10], nil
	}

	return newsLinks, nil
}

// crawlingNewsDetail crawling news detail
func (c *CrawlingNewsCommand) crawlingNewsDetail(newsLink string) (*entity.NewsDetail, error) {
	var (
		collector = c.crawlingSvc.GetCollector()
		news      = &entity.NewsDetail{Link: newsLink}
	)

	// publish time
	collector.OnHTML(valueobject.Html_time, func(h *colly.HTMLElement) {
		publishedStr := h.Attr(valueobject.Attr_datetime)

		data := strings.Split(publishedStr, "T")

		if len(data) < 1 {
			return
		}

		publishedAt, err := time.Parse(time.DateOnly, data[0])
		if err != nil {
			return
		}

		news.PublishedAt = publishedAt
	})

	// title
	collector.OnHTML(valueobject.Html_h1, func(h *colly.HTMLElement) {
		news.Title = h.Text
	})

	// content
	collector.OnHTML(valueobject.Html_p, func(h *colly.HTMLElement) {
		content := strings.TrimSpace(h.Text)

		if len(content) == 0 {
			return
		}

		news.Contents = append(news.Contents, content)
	})

	// images
	collector.OnHTML(valueobject.Html_img, func(h *colly.HTMLElement) {
		imgUrl := strings.TrimSpace(h.Attr(valueobject.Attr_src))

		if len(imgUrl) == 0 {
			return
		}

		news.Images = append(news.Images, urlPrefixHandle(imgUrl, h.Request.URL))
	})

	if err := collector.Visit(newsLink); err != nil && !pkgCollector.IgnorableError(err) {
		return nil, errors.WithStack(err)
	}

	// validate news
	if err := news.Validate(); err != nil {
		return nil, nil
	}

	return news, nil
}

// urlPrefixHandle url prefix handle
func urlPrefixHandle(webUrl string, reqUrl *url.URL) string {
	if strings.HasPrefix(webUrl, "http") {
		return webUrl
	}

	if !strings.Contains(webUrl, reqUrl.Host) {
		webUrl = copyUrl(reqUrl).JoinPath(webUrl).String()
	}

	if strings.HasPrefix(webUrl, "/") {
		return ""
	}

	return webUrl
}

// copyUrl copy url
func copyUrl(src *url.URL) *url.URL {
	return &url.URL{
		Host:   src.Host,
		Scheme: src.Scheme,
	}
}

// isTopicLink is topic link
func isTopicLink(text string, topics []string) bool {
	return slices.ContainsFunc(topics, func(category string) bool {
		return strings.EqualFold(text, category)
	})
}
