package model

import "gorm.io/gorm"

// AutoMigrate will migrate all models to database
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&NewsWebsite{}, &NewsDetail{}, &CrawlingRecord{}, &NewsKeywords{})
}
