package command

import (
	"context"

	"github.com/mjiee/world-news/backend/pkg/audio"
	"github.com/mjiee/world-news/backend/pkg/pathx"
	"github.com/mjiee/world-news/backend/pkg/ttsai"
	"github.com/mjiee/world-news/backend/service"
)

// TextToSpeechCommand represents a command to text to speech.
type TextToSpeechCommand struct {
	batchNo string
	script  *ttsai.TtsScript

	systemConfigSvc service.SystemConfigService
}

func NewTextToSpeechCommand(
	batchNo string,
	script *ttsai.TtsScript,
	systemConfigSvc service.SystemConfigService,
) *TextToSpeechCommand {
	return &TextToSpeechCommand{
		batchNo:         batchNo,
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

	if c.script.Format == "" {
		c.script.Format = audio.WAV
	}

	resp, err := ttsClient.TextToSpeech(ctx, c.script)
	if err != nil {
		return "", err
	}

	audioFile, err := pathx.GetFilePath(resp.AudioId+"."+c.script.Format, pathx.AudioDir, c.batchNo)
	if err != nil {
		return "", err
	}

	err = audio.SaveAudio(resp.AudioData, audioFile)

	return audioFile, err
}
