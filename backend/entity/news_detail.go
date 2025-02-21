package entity

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/urlx"
	"github.com/mjiee/world-news/backend/repository/model"

	"github.com/gocolly/colly/v2"
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
		PublishedAt: m.PublishedAt,
		CreatedAt:   m.CreatedAt,
	}, nil
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
		PublishedAt: n.PublishedAt,
		CreatedAt:   n.CreatedAt,
	}, nil
}

// Validate validates the NewsDetail entity.
func (n *NewsDetail) Validate() error {
	if n == nil {
		return errorx.NewsNotFound
	}

	if n.Title == "" { // title is required
		return errorx.NewsNotFound
	}

	if len(n.Contents) == 0 { // contents is required
		return errorx.NewsNotFound
	}

	return nil
}

// ExtractTitle extracts the title from the news detail page.
func (n *NewsDetail) ExtractTitle(selector string) (string, colly.HTMLCallback) {
	if selector == "" {
		selector = valueobject.Html_h1
	}

	return selector, func(h *colly.HTMLElement) {
		n.Title = h.Text
	}
}

// ExtractPublishTime extracts the publish time from the news detail page.
func (n *NewsDetail) ExtractPublishTime(selector string) (string, colly.HTMLCallback) {
	if selector == "" {
		selector = valueobject.Html_time
	}

	return selector, func(h *colly.HTMLElement) {
		publishedStr := h.Attr(valueobject.Attr_datetime)

		data := strings.Split(publishedStr, "T")

		if len(data) < 1 {
			return
		}

		publishedAt, err := time.Parse(time.DateOnly, data[0])
		if err != nil {
			return
		}

		n.PublishedAt = publishedAt
	}
}

// ExtractContent extracts the content from the news detail page.
func (n *NewsDetail) ExtractContent(selector string) (string, colly.HTMLCallback) {
	if selector == "" {
		selector = valueobject.Html_p
	}

	return selector, func(h *colly.HTMLElement) {
		content := strings.TrimSpace(h.Text)

		if len(content) == 0 {
			return
		}

		n.Contents = append(n.Contents, content)
	}
}

// ExtractImage extracts the image from the news detail page.
func (n *NewsDetail) ExtractImage(selector string) (string, colly.HTMLCallback) {
	if selector == "" {
		selector = valueobject.Html_img
	}

	return selector, func(h *colly.HTMLElement) {
		imgUrl := strings.TrimSpace(h.Attr(valueobject.Attr_src))

		if len(imgUrl) == 0 {
			return
		}

		n.Images = append(n.Images, urlx.UrlPrefixHandle(imgUrl, h.Request.URL))
	}
}
