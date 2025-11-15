package entity

import (
	"time"

	"github.com/mjiee/world-news/backend/entity/valueobject"
)

// Podcast represents the podcast entity.
type Podcast struct {
	Id        uint
	News      *NewsDetail
	Script    string
	Language  string
	Audio     *valueobject.PodcastAudio
	Style     string
	TtsAi     *valueobject.TaskAi
	CreatedAt time.Time
}
