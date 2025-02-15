package valueobject

import "github.com/mjiee/world-news/backend/pkg/httpx"

// QueryRecordParams query record params
type QueryRecordParams struct {
	RecordType string
	Status     string
	Page       *httpx.Pagination
}

// NewQueryRecordParams creates a new QueryRecordParams instance.
func NewQueryRecordParams(recordType, status string, page *httpx.Pagination) *QueryRecordParams {
	return &QueryRecordParams{
		RecordType: recordType,
		Status:     status,
		Page:       page,
	}
}
