package valueobject

import (
	"github.com/mjiee/world-news/backend/pkg/ttsai"
)

// PodcastAudio represents the podcast audio
type PodcastAudio struct {
	Voices   []*ttsai.Voice     `json:"voices"`
	Type     string             `json:"type"` // mp3, wav, m3u8
	Url      string             `json:"url"`
	Data     string             `json:"data"`
	Duration int                `json:"duration"` // s
	Scripts  []*ttsai.TtsScript `json:"scripts"`
}
