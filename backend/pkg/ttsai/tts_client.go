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
	Content    string  `json:"content"`
	Speaker    string  `json:"speaker"`
	Emotion    string  `json:"emotion"`
	SpeechRate float32 `json:"speechRate"` // [0,2], default 1
	Volume     int     `json:"volume"`     // [0,100], default 50
}

// TtsTask represents the response from the tts service
type TtsTask struct {
	AudioId     string
	Type        string // mav, mp3, m3u8
	Url         string
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
