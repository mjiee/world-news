package valueobject

import "github.com/mjiee/world-news/backend/pkg/httpx"

// QueryNewsParams query news params
type QueryNewsParams struct {
	RecordId uint
	Page     *httpx.Pagination
}

// NewQueryNewsParams creates a new QueryNewsParams instance.
func NewQueryNewsParams(recordId uint, page *httpx.Pagination) *QueryNewsParams {
	return &QueryNewsParams{
		RecordId: recordId,
		Page:     page,
	}
}
