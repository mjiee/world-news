package valueobject

import (
	"github.com/mjiee/world-news/backend/pkg/openai"
	"github.com/mjiee/world-news/backend/pkg/ttsai"
)

// TaskAi represents the AI configuration for a task
type TaskAi struct {
	SessionId string `json:"sessionId"`
	Platform  string `json:"platform"`
	Model     string `json:"model"`
}

// NewTaskAiFromTextAi creates a new TaskAi object from TextAiConfig
func NewTaskAiFromTextAi(textAi *openai.Config) *TaskAi {
	return &TaskAi{
		Platform: textAi.Platform,
		Model:    textAi.Model,
	}
}

// NewTaskAiFromTtsAi creates a new TaskAi object from TtsAiConfig
func NewTaskAiFromTtsAi(ttsAi *ttsai.Config) *TaskAi {
	return &TaskAi{
		Platform: ttsAi.Platform,
		Model:    ttsAi.Model,
	}
}
