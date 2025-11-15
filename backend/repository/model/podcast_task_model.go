package model

import "time"

// PodcastTask is the model for podcast task
type PodcastTask struct {
	ID        uint   `gorm:"primaryKey"`
	BatchNo   string `gorm:"index,not null"`
	Stage     string
	Title     string
	NewsId    uint
	Status    string
	Language  string
	Prompt    string
	Input     string
	Output    string
	Reason    string
	Result    string
	Audio     string
	TaskAi    string
	Extra     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
