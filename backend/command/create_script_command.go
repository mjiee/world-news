package command

import (
	"context"
	"slices"
	"strings"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/openai"
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
	textAi, ttsAi, prompt, err := getPodcastConfig(ctx, c.systemConfigSvc)
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

	// execute
	go c.createScript(task, textAi, prompt)

	return nil
}

// createScript create script
func (c *CreateScriptCommand) createScript(task *entity.PodcastTask, textAi *openai.Config,
	prompt *valueobject.PodcastScriptPrompt) {
	var (
		stage    = task.GetStage(valueobject.TaskStageScripted)
		messages = []*openai.Message{
			openai.SystemMessage(prompt.BuildSystemPrompt(task.Language)),
			openai.UserMessage(stage.BuildPrompt()),
		}
	)

	data, err := openai.NewOpenaiClient(textAi).SetMessage(messages...).ChatCompletion(c.ctx)
	if err != nil {
		stage.Fail(err.Error())
		task.Result = valueobject.TaskResultFailed
	} else {
		assistantMsg := gokit.SliceMap(data.Choices, func(item *openai.ChatCompletionChoice) string {
			return item.Message.Content
		})

		stage.TaskAi.SessionId = data.ID
		stage.SetOutput(strings.Join(assistantMsg, "\n"))

		scripts := prompt.ExtractScripts(stage.Output)
		if scripts == nil {
			stage.Fail("failed to extract scripts")
			stage.SetOutput(stage.Output)
		} else {
			stage.Audio.Scripts = scripts
			stage.SetStatus(valueobject.StageStatusCompleted)
		}
	}

	if err := c.taskSvc.SaveTask(c.ctx, task); err != nil {
		logx.WithContext(c.ctx).Error("CreateScriptCommand.SaveTask", err)
	}
}
