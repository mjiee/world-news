package dto

import "github.com/mjiee/world-news/backend/pkg/httpx"

// AddNewsWebsiteRequest add news website request
type AddNewsWebsiteRequest struct {
	WebsiteType string   `json:"websiteType"`
	URL         string   `json:"url"`
	Config      []string `json:"config"`
}

// DeleteNewsWebsiteRequest delete news website request
type DeleteNewsWebsiteRequest struct {
	Id uint `json:"id"` // ID
}

// GetNewsWebsitesRequest get news websites request
type GetNewsWebsitesRequest struct {
	WebsiteType string            `json:"websiteType"`
	Pagination  *httpx.Pagination `json:"pagination"`
}

// GetNewsWebsitesResult get news websites result
type GetNewsWebsitesResult struct {
	Data  []*NewsWebsite `json:"data"`
	Total int64          `json:"total"`
}

// GetNewsWebsitesResponse get news websites response
type GetNewsWebsitesResponse struct {
	*httpx.Response
	Result *GetNewsWebsitesResult `json:"result"`
}

// NewsWebsite news website
type NewsWebsite struct {
	Id          uint     `json:"id"`
	WebsiteType string   `json:"websiteType"`
	URL         string   `json:"url"`
	Config      []string `json:"config"`
}

// GetNewsKeywordsResponse get news keywords response
type GetNewsKeywordsResponse struct {
	*httpx.Response
	Result []string `json:"result"`
}

// DeleteNewsKeywordRequest delete news keywords request
type DeleteNewsKeywordRequest struct {
	Keyword uint `json:"keyword"`
}

// AddNewsKeywordRequest add news keywords request
type AddNewsKeywordRequest struct {
	Keyword string `json:"keyword"`
}
