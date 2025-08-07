package entity

import (
	"encoding/json"
	"time"

	"github.com/mjiee/gokit/slicex"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/textx"
	"github.com/mjiee/world-news/backend/pkg/timex"
	"github.com/mjiee/world-news/backend/pkg/urlx"
	"github.com/mjiee/world-news/backend/repository/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// NewsDetail represents the detailed information about a news item.
type NewsDetail struct {
	Id          uint
	RecordId    uint // crawling record id
	Source      string
	Topic       string
	Title       string
	Author      string
	PublishedAt time.Time
	Link        string
	Contents    []string
	Images      []string
	Video       string
	Scraped     bool
	CreatedAt   time.Time
}

// NewNewsDetailFromModel converts a NewsDetailModel to a NewsDetail entity.
func NewNewsDetailFromModel(m *model.NewsDetail) (*NewsDetail, error) {
	if m == nil {
		return nil, errorx.NewsNotFound
	}

	var (
		contents []string
		images   []string
	)

	if err := json.Unmarshal([]byte(m.Contents), &contents); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := json.Unmarshal([]byte(m.Images), &images); err != nil {
		return nil, errors.WithStack(err)
	}

	return &NewsDetail{
		Id:          m.ID,
		RecordId:    m.RecordId,
		Source:      m.Source,
		Topic:       m.Topic,
		Title:       m.Title,
		Author:      m.Author,
		Link:        m.Link,
		Contents:    contents,
		Images:      images,
		Video:       m.Video,
		Scraped:     m.Scraped,
		PublishedAt: m.PublishedAt,
		CreatedAt:   m.CreatedAt,
	}, nil
}

// NewNewsDetailFromTopicLink creates a NewsDetail entity from a NewsTopicLink.
func NewNewsDetailFromTopicLink(recordId uint, link *valueobject.NewsTopicLink) *NewsDetail {
	return &NewsDetail{
		RecordId: recordId,
		Source:   urlx.ExtractSecondLevelDomain(link.URL),
		Topic:    link.Topic,
	}
}

// ToModel converts the NewsDetail entity to a NewsDetailModel.
func (n *NewsDetail) ToModel() (*model.NewsDetail, error) {
	if n == nil {
		return nil, errorx.NewsNotFound
	}

	contents, err := json.Marshal(n.Contents)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	images, err := json.Marshal(n.Images)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &model.NewsDetail{
		ID:          n.Id,
		RecordId:    n.RecordId,
		Source:      n.Source,
		Topic:       n.Topic,
		Title:       n.Title,
		Author:      n.Author,
		Link:        n.Link,
		Contents:    string(contents),
		Images:      string(images),
		Video:       n.Video,
		Scraped:     n.Scraped,
		PublishedAt: n.PublishedAt,
		CreatedAt:   n.CreatedAt,
	}, nil
}

// IsValid checks if the news detail is valid.
func (n *NewsDetail) IsValid(minPublishTime time.Time) bool {
	if n.PublishedAt.IsZero() {
		return false
	}

	if !minPublishTime.IsZero() && n.PublishedAt.Before(minPublishTime) {
		return false
	}

	return isNewsTitle(n.Title) && n.Link != ""
}

// Compare compare the priorities of news.
func (n *NewsDetail) Compare(other *NewsDetail) int {
	if len(n.Images) > len(other.Images) && len(n.Contents) > len(other.Contents) {
		return 1
	}

	if len(n.Images) < len(other.Images) && len(n.Contents) < len(other.Contents) {
		return -1
	}

	return 0
}

// ExtractTitle extracts the title from the news detail.
func (n *NewsDetail) ExtractTitle(doc *goquery.Selection) {
	for _, selector := range valueobject.NewsTitleSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			if n.Title != "" {
				return
			}

			text := textx.CleanText(s.Text())
			if isNewsTitle(text) {
				n.Title = text
			}
		})

		if n.Title != "" {
			break
		}
	}

	return
}

// ExtractSummary extracts the summary from the news detail.
func (n *NewsDetail) ExtractSummary(doc *goquery.Selection) {
	for _, selector := range valueobject.NewsSummarySelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			if len(n.Contents) > 0 {
				return
			}

			text := textx.CleanText(s.Text())
			if textx.SimilarText(text, n.Title) {
				return
			}

			if isNewsContent(text) {
				n.Contents = append(n.Contents, text)
			}
		})
	}

	if len(n.Contents) > 2 {
		n.Contents = n.Contents[:1]
	}

	return
}

// ExtractNewsDetail extracts the news detail from the document.
func (n *NewsDetail) ExtractNewsDetail(doc *goquery.Selection) {
	var ok bool

	for _, selector := range valueobject.NewsMainBodySelectors {
		bodyDoc := doc.Find(selector)

		ok = n.extractContents(bodyDoc, valueobject.NewsContentSelectors)
		if !ok {
			continue
		}

		n.ExtractAuthor(bodyDoc)

		if len(n.Images) == 0 {
			n.ExtractImages(bodyDoc)
		}

		break
	}

	if !ok {
		_ = n.extractContents(doc, valueobject.NewsContentSelectors)

		n.ExtractAuthor(doc)

		if len(n.Images) == 0 {
			n.ExtractImages(doc)
		}
	}

	if len(n.Images) > 3 {
		if images := slicex.Filter(n.Images, n.isValidImage); len(images) > 0 {
			n.Images = images[:min(3, len(images))]
		} else {
			n.Images = n.Images[:3]
		}
	}
}

// extractContents extracts the contents from the news detail.
func (n *NewsDetail) extractContents(doc *goquery.Selection, selectors []string) bool {
	var contents []string

	for _, selector := range selectors {
		contents = make([]string, 0)

		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			text := textx.CleanText(s.Text())
			if textx.SimilarText(text, n.Title) {
				return
			}

			if isNewsContent(text) {
				contents = append(contents, text)
			}
		})

		if len(contents) > 2 {
			break
		}
	}

	if len(contents) == 0 {
		return false
	}

	n.Contents = contents

	return true
}

// ExtractLink extracts the link from the news detail.
func (n *NewsDetail) ExtractLink(baseURL string, doc *goquery.Selection) {
	doc.Find(valueobject.LinkSelector).Each(func(i int, s *goquery.Selection) {
		if n.Link != "" {
			return
		}

		href, exists := s.Attr(valueobject.Attr_href)
		if !exists {
			return
		}

		if valueobject.IsNewsLink(href) {
			n.Link = urlx.NormalizeURL(baseURL, href)
		}
	})

	if n.Link != "" || !doc.Is(valueobject.Html_a) {
		return
	}

	if href, exists := doc.Attr(valueobject.Attr_href); exists {
		n.Link = urlx.NormalizeURL(baseURL, href)
	}

	return
}

// ExtractImages extracts the images from the news detail.
func (n *NewsDetail) ExtractImages(doc *goquery.Selection) {
	for _, selector := range valueobject.NewsImageSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			n.extractImageLinks(s)
		})
	}

	n.Images = slicex.Distinct(n.Images, func(item string) string { return urlx.RemoveQueryParams(item) })
}

// extractImageLinks extracts the image links from the news detail.
func (n *NewsDetail) extractImageLinks(doc *goquery.Selection) {
	for _, attr := range valueobject.ImageAttrs {
		link, _ := doc.Attr(attr)

		imageUrl := urlx.NormalizeURL(n.Link, link)

		if valueobject.IsNewsImageLink(imageUrl) {
			n.Images = append(n.Images, imageUrl)
		}
	}
}

// ExtractPublishTime extracts the publish time from the news detail.
func (n *NewsDetail) ExtractPublishTime(doc *goquery.Selection) {
	// try to extract publish time from selector
	for _, selector := range valueobject.NewsTimeSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			if !n.PublishedAt.IsZero() {
				return
			}

			n.PublishedAt = getTimeFromElement(s)
		})
	}

	// try to extract publish time from text
	doc.Each(func(i int, s *goquery.Selection) {
		if !n.PublishedAt.IsZero() {
			return
		}

		n.PublishedAt = timex.ParseTime(s.Text())
	})

	if !n.PublishedAt.IsZero() {
		return
	}

	// try to extract publish time from link
	if n.Link != "" {
		n.PublishedAt = timex.ParseTime(n.Link)
	}

	if !n.PublishedAt.IsZero() {
		return
	}

	if n.Video != "" {
		n.PublishedAt = timex.ParseTime(n.Video)
	}

	if !n.PublishedAt.IsZero() {
		return
	}

	for _, link := range n.Images {
		n.PublishedAt = timex.ParseTime(link)
	}
}

// getTimeFromElement get time from element
func getTimeFromElement(s *goquery.Selection) time.Time {
	for _, attr := range valueobject.TimeAttributes {
		val, _ := s.Attr(attr)

		date := timex.ParseTime(val)
		if !date.IsZero() {
			return date
		}
	}

	return timex.ParseTime(s.Text())
}

// ExtractAuthor extracts the author from the news detail.
func (n *NewsDetail) ExtractAuthor(doc *goquery.Selection) {
	for _, selector := range valueobject.NewsAuthorSelectors {
		if n.Author != "" {
			return
		}

		doc.Filter(selector).Each(func(i int, s *goquery.Selection) {
			if n.Author != "" {
				return
			}

			text := textx.CleanText(s.Text())
			if s.Is(valueobject.Html_meta) {
				text, _ = s.Attr(valueobject.Attr_content)
			}

			if isAuthor(text) {
				n.Author = text
			}
		})
	}
}

// isValidImage checks if a given image is valid.
func (n *NewsDetail) isValidImage(imageLink string) bool {
	if imageLink == "" {
		return false
	}

	if n.PublishedAt.IsZero() {
		return true
	}

	var (
		imageTime   = timex.ParseTime(imageLink)
		publishTime = n.PublishedAt.Add(-(valueobject.MaxValidityPeriod))
	)

	return !imageTime.IsZero() && imageTime.After(publishTime)
}

// isNewsTitle checks if a given news title is valid.
func isNewsTitle(title string) bool {
	return len(title) >= 20 && len(title) <= 500
}

// isNewsContent checks if a given news content is valid.
func isNewsContent(content string) bool {
	return len(content) >= 50
}

// isAuthor checks if a given news author is valid.
func isAuthor(author string) bool {
	return author != "" && len(author) >= 1 && len(author) <= 100
}
