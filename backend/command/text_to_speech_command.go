package command

import (
	"context"
	"encoding/base64"

	"github.com/mjiee/world-news/backend/pkg/ttsai"
	"github.com/mjiee/world-news/backend/service"
)

// TextToSpeechCommand represents a command to text to speech.
type TextToSpeechCommand struct {
	script *ttsai.TtsScript

	systemConfigSvc service.SystemConfigService
}

func NewTextToSpeechCommand(
	script *ttsai.TtsScript,
	systemConfigSvc service.SystemConfigService,
) *TextToSpeechCommand {
	return &TextToSpeechCommand{
		script:          script,
		systemConfigSvc: systemConfigSvc,
	}
}

func (c *TextToSpeechCommand) Execute(ctx context.Context) (string, error) {
	// get config
	_, ttsAi, _, err := c.systemConfigSvc.GetPodcastConfig(ctx)
	if err != nil {
		return "", nil
	}

	ttsClient, err := ttsai.NewDoubaoTTSClient(ttsAi)
	if err != nil {
		return "", err
	}

	resp, err := ttsClient.TextToSpeech(ctx, c.script)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(resp.AudioData), nil
}
