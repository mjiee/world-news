package model

import "time"

// NewsDetail represents the detailed information about a news item.
type NewsDetail struct {
	ID          uint `gorm:"primaryKey"`
	RecordId    uint `gorm:"index;not null"` // crawling record id
	Title       string
	Link        string
	Contents    string
	Images      string
	PublishedAt time.Time
	CreatedAt   time.Time
}

func (n *NewsDetail) TableName() string {
	return "news_details"
}
