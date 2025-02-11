package valueobject

// CrawlingRecordStatus represents the status of a crawling record.
type CrawlingRecordStatus string

const (
	ProcessingCrawlingRecord CrawlingRecordStatus = "processing"
	CompletedCrawlingRecord  CrawlingRecordStatus = "completed"
	FailedCrawlingRecord     CrawlingRecordStatus = "failed"
)
