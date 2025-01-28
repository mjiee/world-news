package dto

import "net/http"

// CrawlingNewsRequest is a struct for requesting news crawling tasks.
type CrawlingNewsRequest struct {
	StartTime string `json:"startTime"`
}

// GetCrawlingRecordsRequest is a struct for requesting crawling records.
type GetCrawlingRecordsRequest struct {
	Page  uint `json:"page,omitempty"`
	Limit uint `json:"limit,omitempty"`
}

// GetCrawlingRecordResult is the result struct for crawling records.
type GetCrawlingRecordResult struct {
	Data  []*CrawlingRecord `json:"data"`
	Total int64             `json:"total"`
}

// GetCrawlingRecordsResponse is the response struct for crawling records.
type GetCrawlingRecordsResponse struct {
	*http.Response
	Result *GetCrawlingRecordResult `json:"result"`
}

// CrawlingRecord represents a single crawling record.
type CrawlingRecord struct {
	Id       uint   `json:"id"`
	Date     string `json:"date"`
	Quantity int64  `json:"quantity"`
	Status   string `json:"status"`
}

// DeleteCrawlingRecordRequest is a struct for deleting a crawling record.
type DeleteCrawlingRecordRequest struct {
	Id uint `json:"id"`
}
