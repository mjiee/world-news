package entity

import "time"

// NewsWebsite represents the basic information of a news website.
type NewsWebsite struct {
	Id          uint `gorm:"primaryKey"`
	URL         string
	WebsiteType WebsiteType
	CreatedAt   time.Time
}

func (n *NewsWebsite) TableName() string {
	return "news_websites"
}

// WebsiteType represents the type of a news website.
type WebsiteType string

const (
	AggregationWebsiteType WebsiteType = "AggregationWebsite"
	NewsWebsiteType        WebsiteType = "NewsWebsite"
)
