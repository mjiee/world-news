package ttsai

import (
	"context"
)

type TtsClient interface {
	// Voices gets the speakers
	Voices(context.Context, string) ([]*Voice, error)
	// TextToSpeech converts text to speech
	TextToSpeech(context.Context, *CreateTtsRequest) (*TtsTask, error)
	// GetAudio gets the audio
	GetAudio(context.Context, string) (*TtsTask, error)
}

// Config for tts service
type Config struct {
	Platform string   `json:"platform"`
	AppId    string   `json:"appId"`
	ApiKey   string   `json:"apiKey"`
	Model    string   `json:"model"`
	Voices   []*Voice `json:"voices"`
}

// Voice represents a speaker
type Voice struct {
	Id          string `json:"id"`                    // voice id
	Name        string `json:"name"`                  // speaker name
	Role        string `json:"role,omitempty"`        // role
	Description string `json:"description,omitempty"` // description
	Model       string `json:"model,omitempty"`
}

// CreateTtsRequest to convert text to speech
type CreateTtsRequest struct {
	Scripts  []*TtsScript
	Language string
}

// TtsScript represents the tts script
type TtsScript struct {
	Text     string  `json:"text"`
	Format   string  `json:"format"` // mp3, mav
	Speaker  string  `json:"speaker"`
	Emotion  string  `json:"emotion"`
	Speed    float32 `json:"speed"`              // [0,2], default 1
	Volume   int     `json:"volume"`             // [0,100], default 50
	Silence  float32 `json:"silence"`            // [0.0, 5.0], default 0.2
	AudioUrl string  `json:"audioUrl,omitempty"` // audio url
}

// TtsTask represents the response from the tts service
type TtsTask struct {
	AudioId     string
	Format      string // mav, mp3, m3u8
	AudioUrl    string
	AudioData   []byte
	Script      string
	Status      ProcessStatus // status: processing, completed, failed
	FailMessage string        // fail message
}

// ProcessStatus represents the status of the process
type ProcessStatus string

const (
	ProcessStatusProcessing ProcessStatus = "processing"
	ProcessStatusCompleted  ProcessStatus = "completed"
	ProcessStatusFailed     ProcessStatus = "failed"
)
