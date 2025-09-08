package valueobject

import (
	"time"

	"github.com/mjiee/world-news/backend/pkg/httpx"
)

// QueryNewsParams query news params
type QueryNewsParams struct {
	RecordId    uint
	Source      string
	Topic       string
	PublishDate time.Time
	Favorited   bool
	Page        *httpx.Pagination
}

// NewQueryNewsParams creates a new QueryNewsParams instance.
func NewQueryNewsParams(recordId uint, page *httpx.Pagination) *QueryNewsParams {
	return &QueryNewsParams{
		RecordId: recordId,
		Page:     page,
	}
}
