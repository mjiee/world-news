package valueobject

import (
	"time"

	"github.com/mjiee/world-news/backend/pkg/httpx"
)

// QueryPodcastTaskParams query podcast task params
type QueryPodcastTaskParams struct {
	StartDate time.Time
	EndDate   time.Time
	Page      *httpx.Pagination
}
