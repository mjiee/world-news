package command

import (
	"context"
	"encoding/base64"
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
	_, ttsAi, _, err := c.systemConfigSvc.GetPodcastConfig(ctx)
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

	stage.Audio = &valueobject.PodcastAudio{Voices: voices}
	task.AddNewStage(stage)

	if err := c.taskSvc.SaveTask(c.ctx, task); err != nil {
		return err
	}

	go c.textToSpeech(task, ttsAi)

	return nil
}

func (c *CreateAudioCommand) textToSpeech(task *entity.PodcastTask, ttsAi *ttsai.Config) error {
	ttsClient, err := ttsai.NewDoubaoTTSClient(ttsAi)
	if err != nil {
		return err
	}

	var (
		stage       = task.GetStage(valueobject.TaskStageTextToSpeech)
		scriptState = task.GetStage(valueobject.TaskStageScripted)
	)

	for _, script := range scriptState.Audio.Scripts {
		if script.Text == "" {
			continue
		}

		_, err = c.taskSvc.GetTaskStage(c.ctx, stage.Id)
		if err != nil {
			stage.Fail(err.Error())
			task.Result = valueobject.TaskResultFailed
			logx.WithContext(c.ctx).Error("GeneratePodcastCommand.textToSpeech", err)
			break
		}

		script.Format = audio.MP3

		resp, err := ttsClient.TextToSpeech(c.ctx, script)
		if err != nil {
			stage.Fail(err.Error())
			task.Result = valueobject.TaskResultFailed
			logx.WithContext(c.ctx).Error("GeneratePodcastCommand.textToSpeech", err)
			break
		}

		stage.Audio.Format = resp.Format
		if len(resp.AudioData) != 0 {
			script.Audio = base64.StdEncoding.EncodeToString(resp.AudioData)
		}
	}

	if !task.Result.IsFailed() {
		task.Result = valueobject.TaskResultCompleted
		stage.SetStatus(valueobject.StageStatusCompleted)
	}

	if err := c.taskSvc.SaveTask(c.ctx, task); err != nil {
		logx.WithContext(c.ctx).Error("GeneratePodcastCommand.textToSpeech", err)
	}

	return err
}
