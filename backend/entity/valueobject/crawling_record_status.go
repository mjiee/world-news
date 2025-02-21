package valueobject

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
