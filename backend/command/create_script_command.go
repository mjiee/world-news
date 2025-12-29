package command

import (
	"context"
	"slices"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/ttsai"
	"github.com/mjiee/world-news/backend/service"
)

// CreateScriptCommand is a command to create script
type CreateScriptCommand struct {
	ctx      context.Context
	stageId  uint
	voiceIds []string

	systemConfigSvc service.SystemConfigService
	taskSvc         service.PodcastTaskService
}

func NewCreateScriptCommand(
	ctx context.Context,
	stageId uint,
	voiceIds []string,
	systemConfigSvc service.SystemConfigService,
	taskSvc service.PodcastTaskService,
) *CreateScriptCommand {
	return &CreateScriptCommand{
		ctx:             ctx,
		stageId:         stageId,
		voiceIds:        voiceIds,
		systemConfigSvc: systemConfigSvc,
		taskSvc:         taskSvc,
	}
}

func (c *CreateScriptCommand) Execute(ctx context.Context) error {
	// config
	textAi, ttsAi, _, err := c.systemConfigSvc.GetPodcastConfig(ctx)
	if err != nil {
		return err
	}

	task, err := c.taskSvc.GetTaskByStageId(ctx, c.stageId)
	if err != nil {
		return err
	}

	if err = task.VerifyTask(); err != nil {
		return err
	}

	voices := gokit.SliceFilter(ttsAi.Voices, func(v *ttsai.Voice) bool { return slices.Contains(c.voiceIds, v.Id) })
	if len(voices) != len(c.voiceIds) {
		return errorx.PodcastVoiceNotFound
	} else if len(c.voiceIds) == 0 {
		if len(ttsAi.Voices) == 0 {
			return errorx.PodcastVoiceNotFound
		}

		voices = []*ttsai.Voice{voices[0]}
	}

	// new stage
	var (
		stage = valueobject.NewTaskStage(valueobject.TaskStageScripted, valueobject.BuildScriptPrompt(task.Language, voices),
			valueobject.NewTaskAiFromTextAi(textAi))
		stylizeStage = task.GetStageById(c.stageId)
	)

	stage.Input = stylizeStage.Output
	stage.Audio = &valueobject.PodcastAudio{Voices: voices}
	task.AddNewStage(stage)

	if err = c.taskSvc.SaveTask(ctx, task); err != nil {
		return err
	}

	// execute task
	executeCmd := NewExecuteTaskCommand(c.ctx, task, c.systemConfigSvc, c.taskSvc)
	go func() {
		if err := executeCmd.Execute(executeCmd.ctx); err != nil {
			logx.WithContext(executeCmd.ctx).Error("CreateScriptCommand", err)
		}
	}()

	return nil
}
