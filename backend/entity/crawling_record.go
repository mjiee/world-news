package entity

import "time"

// CrawlingRecord represents a crawling record.
type CrawlingRecord struct {
	Id        uint `gorm:"primaryKey"`
	Date      string
	Quantity  int64
	Status    string
	CreatedAt time.Time
}
