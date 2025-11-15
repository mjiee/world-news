package command

import (
	"context"
	"slices"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/audio"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/ttsai"
	"github.com/mjiee/world-news/backend/service"
)

// CreateAudioCommand represents the command to generate podcast.
type CreateAudioCommand struct {
	ctx     context.Context
	stageId uint

	systemConfigSvc service.SystemConfigService
	taskSvc         service.PodcastTaskService
}

func NewCreateAudioCommand(
	ctx context.Context,
	stageId uint,
	systemConfigSvc service.SystemConfigService,
	taskSvc service.PodcastTaskService,
) *CreateAudioCommand {
	return &CreateAudioCommand{
		ctx:             ctx,
		stageId:         stageId,
		systemConfigSvc: systemConfigSvc,
		taskSvc:         taskSvc,
	}
}

func (c *CreateAudioCommand) Execute(ctx context.Context) error {
	// get config
	_, ttsAi, _, err := getPodcastConfig(ctx, c.systemConfigSvc)
	if err != nil {
		return err
	}

	// get task
	task, err := c.taskSvc.GetTaskByStageId(ctx, c.stageId)
	if err != nil {
		return err
	}

	if err = task.VerifyTask(); err != nil {
		return err
	}

	// new stage
	var (
		stage    = valueobject.NewTaskStage(valueobject.TaskStageTextToSpeech, "", valueobject.NewTaskAiFromTtsAi(ttsAi))
		scripts  = task.GetPodcastScript()
		spickers = gokit.SliceMap(scripts, func(s *ttsai.TtsScript) string { return s.Speaker })
		voices   = gokit.SliceFilter(ttsAi.Voices, func(v *ttsai.Voice) bool { return slices.Contains(spickers, v.Id) })
	)

	if len(scripts) == 0 {
		return errorx.PodcastScriptNotFound
	}

	stage.Audio = &valueobject.PodcastAudio{Scripts: scripts, Voices: voices}
	task.AddNewStage(stage)

	if err := c.taskSvc.SaveTask(c.ctx, task); err != nil {
		return err
	}

	go c.textToSpeech(task, ttsAi)

	return nil
}

func (c *CreateAudioCommand) textToSpeech(task *entity.PodcastTask, ttsAi *ttsai.Config) {
	stage := task.GetStage(valueobject.TaskStageTextToSpeech)

	// create tts
	if err := c.createTts(stage, ttsAi); err != nil {
		stage.Fail(err.Error())
		task.Result = valueobject.TaskResultFailed
	} else {
		stage.Status = valueobject.StageStatusCompleted
		task.Result = valueobject.TaskResultCompleted
	}

	if err := c.taskSvc.SaveTask(c.ctx, task); err != nil {
		logx.WithContext(c.ctx).Error("GeneratePodcastCommand.textToSpeech", err)
	}
}

func (c *CreateAudioCommand) createTts(stage *valueobject.TaskStage, ttsAi *ttsai.Config) error {
	ttsClient, err := ttsai.NewDoubaoTTSClient(ttsAi)
	if err != nil {
		return err
	}

	var (
		sessionId []string
		audioData [][]byte
	)

	for _, script := range stage.Audio.Scripts {
		if script.Content == "" {
			continue
		}

		resp, err := ttsClient.TextToSpeech(c.ctx, script)
		if err != nil {
			return err
		}

		sessionId = append(sessionId, resp.AudioId)
		stage.Audio.Type = resp.Type
		if len(resp.AudioData) != 0 {
			audioData = append(audioData, resp.AudioData)
		}
	}

	if len(audioData) > 0 {
		switch stage.Audio.Type {
		case audio.MP3:
			data, duration, err := audio.EncodeMp3s(audioData...)
			if err != nil {
				return err
			}

			stage.Audio.Data = data
			stage.Audio.Duration = int(duration.Seconds())
		case audio.WAV:
			stage.Audio.Data, err = audio.EncodeWavs(audioData...)
		}
	}

	return err
}
