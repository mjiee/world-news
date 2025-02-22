package valueobject

import "github.com/mjiee/world-news/backend/pkg/errorx"

// CrawlingRecordStatus represents the status of a crawling record.
type CrawlingRecordStatus string

const (
	ProcessingCrawlingRecord CrawlingRecordStatus = "processing"
	CompletedCrawlingRecord  CrawlingRecordStatus = "completed"
	FailedCrawlingRecord     CrawlingRecordStatus = "failed"
	PausedCrawlingRecord     CrawlingRecordStatus = "paused"
)

func (s CrawlingRecordStatus) String() string {
	return string(s)
}

// IsFailed returns true if the crawling record status is failed.
func (s CrawlingRecordStatus) IsFailed() bool {
	return s == FailedCrawlingRecord
}

// IsProcessing returns true if the crawling record status is processing.
func (s CrawlingRecordStatus) IsProcessing() bool {
	return s == ProcessingCrawlingRecord
}

// IsPaused returns true if the crawling record status is paused.
func (s CrawlingRecordStatus) IsPaused() bool {
	return s == PausedCrawlingRecord
}

// UpdateValidStatus updates the crawling record status to valid status.
func (s CrawlingRecordStatus) UpdateValidStatus(newStatus CrawlingRecordStatus) error {
	switch newStatus {
	case PausedCrawlingRecord:
		if !s.IsProcessing() {
			return errorx.UpdateRecordStatusNotAllowed
		}
	case ProcessingCrawlingRecord:
		if !s.IsPaused() {
			return errorx.UpdateRecordStatusNotAllowed
		}
	default:
		return errorx.UpdateRecordStatusNotAllowed
	}

	return nil
}
