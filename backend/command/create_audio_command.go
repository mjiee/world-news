package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/audio"
	"github.com/mjiee/world-news/backend/pkg/config"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/pathx"
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
	var (
		stage       = task.GetStage(valueobject.TaskStageTextToSpeech)
		scriptState = task.GetStage(valueobject.TaskStageScripted)
	)

	audioPath, err := pathx.GetAppBasePath(config.AppName, pathx.AudioDir, task.BatchNo)
	if err != nil {
		return err
	}

	err = c.generateAudio(audioPath, stage, scriptState, ttsAi)
	if err == nil {
		err = c.mergeAudio(audioPath, stage, scriptState)
	}

	if err != nil {
		stage.Fail(err.Error())
		task.Result = valueobject.TaskResultFailed
		logx.WithContext(c.ctx).Error("GeneratePodcastCommand", err)
	}

	if !task.Result.IsFailed() {
		task.Result = valueobject.TaskResultCompleted
		stage.SetStatus(valueobject.StageStatusCompleted)
	}

	if err := c.taskSvc.SaveTask(c.ctx, task); err != nil {
		logx.WithContext(c.ctx).Error("GeneratePodcastCommand.saveTask", err)
	}

	return err
}

func (c *CreateAudioCommand) generateAudio(audioPath string, ttsStage *valueobject.TaskStage,
	scriptState *valueobject.TaskStage, ttsAi *ttsai.Config) error {
	ttsClient, err := ttsai.NewDoubaoTTSClient(ttsAi)
	if err != nil {
		return err
	}

	for _, script := range scriptState.Audio.Scripts {
		if script.Text == "" {
			continue
		}

		_, err := c.taskSvc.GetTaskStage(c.ctx, ttsStage.Id)
		if err != nil {
			return err
		}

		if script.AudioUrl != "" {
			continue
		}

		script.Format = audio.WAV

		resp, err := ttsClient.TextToSpeech(c.ctx, script)
		if err != nil {
			return err
		}

		ttsStage.Audio.Format = resp.Format
		if len(resp.AudioData) == 0 {
			continue
		}

		audioFile := filepath.Join(audioPath, resp.AudioId+"."+script.Format)

		if err := audio.Transcode(resp.AudioData, audioFile); err != nil {
			return err
		}

		script.AudioUrl = audioFile
	}

	return nil
}

func (c *CreateAudioCommand) mergeAudio(audioPath string, ttsStage *valueobject.TaskStage,
	scriptState *valueobject.TaskStage) error {
	var (
		tempAudios   = make([]string, 0, len(scriptState.Audio.Scripts))
		leftSpeacker = scriptState.Audio.Scripts[0].Speaker
		audioFile    = filepath.Join(audioPath, fmt.Sprintf("%s_%d.wav", ttsStage.BatchNo, ttsStage.Id))
	)

	tempPath, err := pathx.GetAppBasePath(config.AppName, pathx.TempDir, ttsStage.BatchNo)
	if err != nil {
		return err
	}

	for _, script := range scriptState.Audio.Scripts {
		if script.AudioUrl == "" {
			continue
		}

		panning := audio.LeftPanning
		if script.Speaker != leftSpeacker {
			panning = audio.RightPanning
		}

		tempFile := filepath.Join(tempPath, filepath.Base(script.AudioUrl))

		if err := audio.RenderFile(script.AudioUrl, tempFile, audio.RenderOption{Pan: panning}); err != nil {
			return err
		}

		tempAudios = append(tempAudios, tempFile)
	}

	if err := audio.MergeFiles(tempAudios, audioFile); err != nil {
		return err
	}

	ttsStage.Audio.Url = audioFile

	if err := os.RemoveAll(tempPath); err != nil {
		logx.WithContext(c.ctx).Error("GeneratePodcastCommand.RemoveAll", err)
	}

	return nil
}
