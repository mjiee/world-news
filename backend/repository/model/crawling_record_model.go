package model

import "time"

// CrawlingRecord represents a crawling record.
type CrawlingRecord struct {
	Id        uint `gorm:"primaryKey"`
	Date      string
	Quantity  int64
	Status    string
	CreatedAt time.Time
}

func (c *CrawlingRecord) TableName() string {
	return "crawling_records"
}
