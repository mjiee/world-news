package valueobject

// CrawlingRecordStatus represents the status of a crawling record.
type CrawlingRecordStatus string

const (
	ProcessingCrawlingRecord CrawlingRecordStatus = "processing"
	CompletedCrawlingRecord  CrawlingRecordStatus = "completed"
	FailedCrawlingRecord     CrawlingRecordStatus = "failed"
)

// IsFailed returns true if the crawling record status is failed.
func (s CrawlingRecordStatus) IsFailed() bool {
	return s == FailedCrawlingRecord
}
