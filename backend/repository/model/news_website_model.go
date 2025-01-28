package model

import "time"

// NewsWebsite represents the basic information of a news website.
type NewsWebsite struct {
	Id          uint   `gorm:"primaryKey"`
	URL         string `gorm:"type:varchar(255);NOT NULL"`
	WebsiteType string `gorm:"type:varchar(40);NOT NULL"`
	Config      string
	CreatedAt   time.Time
}

func (n *NewsWebsite) TableName() string {
	return "news_websites"
}
