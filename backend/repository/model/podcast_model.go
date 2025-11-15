package model

import "time"

// Podcast is the model for podcast
type Podcast struct {
	ID        uint `gorm:"primaryKey"`
	NewsId    uint `gorm:"index,not null"`
	Script    string
	Language  string
	Audio     string
	Style     string
	TtsAi     string
	CreatedAt time.Time
}
