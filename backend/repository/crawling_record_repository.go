package repository

import "gorm.io/gorm"

// CrawlingRecordRepository is interface for crawling record.
type CrawlingRecordRepository interface {
}

type crawlingRecordRepository struct {
	db *gorm.DB
}

func NewCrawlingRecordRepository(db *gorm.DB) CrawlingRecordRepository {
	return &crawlingRecordRepository{db}
}
