package dto

import (
	"time"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/httpx"
)

// CrawlingNewsRequest is a struct for requesting news crawling tasks.
type CrawlingNewsRequest struct {
	StartTime string `json:"startTime,omitempty"`
}

// GetCrawlingRecordRequest is a struct for requesting a crawling record.
type GetCrawlingRecordRequest struct {
	Id uint `json:"id" binding:"required"`
}

// QueryCrawlingRecordsRequest is a struct for requesting crawling records.
type QueryCrawlingRecordsRequest struct {
	RecordType string            `json:"recordType,omitempty"`
	Status     string            `json:"status,omitempty"`
	Pagination *httpx.Pagination `json:"pagination"`
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
		records[i] = mappingCrawlingRecord(record)
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
	Id         uint                  `json:"id"`
	RecordType string                `json:"recordType"`
	Quantity   int64                 `json:"quantity"`
	Status     string                `json:"status"`
	Config     *CrawlingRecordConfig `json:"config,omitempty"`
	StartTime  string                `json:"startTime"`
	EndTime    string                `json:"endTime"`
}

// NewCrawlingRecordFromEntity creates a new CrawlingRecord instance.
func NewCrawlingRecordFromEntity(record *entity.CrawlingRecord) *CrawlingRecord {
	if record == nil {
		return nil
	}

	data := mappingCrawlingRecord(record)
	data.Config = NewCrawlingRecordConfigFromValue(record.Config)

	return data
}

// mappingCrawlingRecord maps a CrawlingRecord entity to a CrawlingRecord struct.
func mappingCrawlingRecord(record *entity.CrawlingRecord) *CrawlingRecord {
	if record == nil {
		return nil
	}

	return &CrawlingRecord{
		Id:         record.Id,
		RecordType: string(record.RecordType),
		Quantity:   record.Quantity,
		Status:     string(record.Status),
		StartTime:  record.CreatedAt.Format(time.DateTime),
		EndTime:    record.UpdatedAt.Format(time.DateTime),
	}
}

// CrawlingRecordConfig represents the configuration of a crawling record.
type CrawlingRecordConfig struct {
	Sources []string `json:"sources,omitempty"`
	Topics  []string `json:"topics,omitempty"`
}

// NewCrawlingRecordConfigFromValue creates a new CrawlingRecordConfig instance.
func NewCrawlingRecordConfigFromValue(data *valueobject.CrawlingRecordConfig) *CrawlingRecordConfig {
	if data == nil {
		return nil
	}

	config := &CrawlingRecordConfig{Topics: data.Topics}

	for _, source := range data.Sources {
		config.Sources = append(config.Sources, source.GetHost())
	}

	return config
}

// DeleteCrawlingRecordRequest is a struct for deleting a crawling record.
type DeleteCrawlingRecordRequest struct {
	Id uint `json:"id" binding:"required"`
}

// UpdateCrawlingRecordStatusRequest is a struct for updating a crawling record.
type UpdateCrawlingRecordStatusRequest struct {
	Id     uint   `json:"id" binding:"required"`
	Status string `json:"status" binding:"oneof=processing paused"`
}
