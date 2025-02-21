package model

import "time"

// CrawlingRecord represents a crawling record.
type CrawlingRecord struct {
	ID         uint `gorm:"primaryKey"`
	RecordType string
	Date       time.Time
	Quantity   int64
	Status     string
	Config     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (c *CrawlingRecord) TableName() string {
	return "crawling_records"
}
