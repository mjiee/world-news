package dto

import (
	"time"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/httpx"
)

// QueryNewsRequest get news detail list request
type QueryNewsRequest struct {
	RecordId   uint              `json:"recordId"`
	Source     string            `json:"source"`
	Topic      string            `json:"topic"`
	Pagination *httpx.Pagination `json:"pagination"`
}

// ToValueobject query news params
func (q *QueryNewsRequest) ToValueobject() *valueobject.QueryNewsParams {
	return &valueobject.QueryNewsParams{
		RecordId: q.RecordId,
		Source:   q.Source,
		Topic:    q.Topic,
		Page:     q.Pagination,
	}
}

// QueryNewsResult get news detail list result
type QueryNewsResult struct {
	Data  []*NewsDetail `json:"data"`
	Total int64         `json:"total"`
}

// NewQueryNewsResult news detail list result
func NewQueryNewsResult(data []*entity.NewsDetail, total int64) *QueryNewsResult {
	res := make([]*NewsDetail, len(data))

	for i, v := range data {
		item := NewNewsDetailFromEntity(v)

		if len(item.Images) > 1 {
			item.Images = item.Images[:1]
		}

		if len(item.Contents) > 1 {
			item.Contents = item.Contents[:1]
		}

		res[i] = item
	}

	return &QueryNewsResult{
		Data:  res,
		Total: total,
	}
}

// QueryNewsResponse get news detail list response
type QueryNewsResponse struct {
	*httpx.Response
	Result *QueryNewsResult `json:"result"`
}

// NewsDetail news detail
type NewsDetail struct {
	Id          uint     `json:"id"`
	Title       string   `json:"title"`
	Source      string   `json:"source"`
	Topic       string   `json:"topic"`
	Link        string   `json:"link"`
	Contents    []string `json:"contents"`
	Images      []string `json:"images"`
	PublishedAt string   `json:"publishedAt"`
}

// NewNewsDetailFromEntity news detail
func NewNewsDetailFromEntity(data *entity.NewsDetail) *NewsDetail {
	if data == nil {
		return nil
	}

	publishedAt := ""

	if !data.PublishedAt.IsZero() {
		publishedAt = data.PublishedAt.Format(time.DateTime)
	}

	return &NewsDetail{
		Id:          data.Id,
		Title:       data.Title,
		Source:      data.Source,
		Topic:       data.Topic,
		Link:        data.Link,
		Contents:    data.Contents,
		Images:      data.Images,
		PublishedAt: publishedAt,
	}
}

// GetNewsDetailRequest get news detail request
type GetNewsDetailRequest struct {
	Id uint `json:"id"`
}

// GetNewsDetailResponse get news detail response
type GetNewsDetailResponse struct {
	*httpx.Response
	Result *NewsDetail `json:"result"`
}

// DeleteNewsRequest delete news detail request
type DeleteNewsRequest struct {
	Id uint `json:"id"`
}
