package dto

import "github.com/mjiee/world-news/backend/pkg/httpx"

// GetNewsListRequest get news detail list request
type GetNewsListRequest struct {
	RecordId   uint              `json:"recordId"`
	Pagination *httpx.Pagination `json:"pagination"`
}

// GetNewsListResult get news detail list result
type GetNewsListResult struct {
	Data  []*NewsDetail `json:"data"`
	Total int64         `json:"total"`
}

// GetNewsListResponse get news detail list response
type GetNewsListResponse struct {
	*httpx.Response
	Result *GetNewsListResult `json:"result"`
}

// NewsDetail news detail
type NewsDetail struct {
	Id          uint     `json:"id"`
	Title       string   `json:"title"`
	Link        string   `json:"link"`
	Contents    []string `json:"contents"`
	Images      []string `json:"images"`
	PublishedAt string   `json:"publishedAt"`
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
