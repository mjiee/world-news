package model

import "time"

// NewsDetail represents the detailed information about a news item.
type NewsDetail struct {
	ID          uint   `gorm:"primaryKey"`
	RecordId    uint   `gorm:"index;not null"` // crawling record id
	Source      string `gorm:"index;not null"`
	Topic       string `gorm:"index;not null"`
	Title       string
	Author      string
	PublishedAt time.Time
	Link        string
	Contents    string
	Images      string
	Video       string
	CreatedAt   time.Time
}

func (n *NewsDetail) TableName() string {
	return "news_details"
}
