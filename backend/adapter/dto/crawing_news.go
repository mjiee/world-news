package dto

import (
	"time"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/pkg/httpx"
)

// CrawlingNewsRequest is a struct for requesting news crawling tasks.
type CrawlingNewsRequest struct {
	StartTime string `json:"startTime"`
}

// QueryCrawlingRecordsRequest is a struct for requesting crawling records.
type QueryCrawlingRecordsRequest struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

// QueryCrawlingRecordResult is the result struct for crawling records.
type QueryCrawlingRecordResult struct {
	Data  []*CrawlingRecord `json:"data"`
	Total int64             `json:"total"`
}

// NewQueryCrawlingRecordResult creates a new GetCrawlingRecordResult instance.
func NewQueryCrawlingRecordResult(data []*entity.CrawlingRecord, total int64) *QueryCrawlingRecordResult {
	records := make([]*CrawlingRecord, len(data))

	for i, record := range data {
		records[i] = NewCrawlingRecordFromEntity(record)
	}

	return &QueryCrawlingRecordResult{
		Data:  records,
		Total: total,
	}
}

// QueryCrawlingRecordsResponse is the response struct for crawling records.
type QueryCrawlingRecordsResponse struct {
	*httpx.Response
	Result *QueryCrawlingRecordResult `json:"result"`
}

// CrawlingRecord represents a single crawling record.
type CrawlingRecord struct {
	Id         uint   `json:"id"`
	RecordType string `json:"recordType"`
	Date       string `json:"date"`
	Quantity   int64  `json:"quantity"`
	Status     string `json:"status"`
}

// NewCrawlingRecordFromEntity creates a new CrawlingRecord instance.
func NewCrawlingRecordFromEntity(record *entity.CrawlingRecord) *CrawlingRecord {
	return &CrawlingRecord{
		Id:         record.Id,
		RecordType: string(record.RecordType),
		Date:       record.Date.Format(time.DateTime),
		Quantity:   record.Quantity,
		Status:     string(record.Status),
	}
}

// DeleteCrawlingRecordRequest is a struct for deleting a crawling record.
type DeleteCrawlingRecordRequest struct {
	Id uint `json:"id"`
}
