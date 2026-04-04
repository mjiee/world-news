package command

import (
	"context"
	"io"
	"strings"

	"github.com/cloudwego/eino/schema"
	"github.com/pkg/errors"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/openai"
	"github.com/mjiee/world-news/backend/service"
)

// ExecuteTaskCommand represents the command to execute a task
type ExecuteTaskCommand struct {
	ctx             context.Context
	task            *entity.PodcastTask
	systemConfigSvc service.SystemConfigService
	taskSvc         service.PodcastTaskService
}

func NewExecuteTaskCommand(
	ctx context.Context,
	task *entity.PodcastTask,
	systemConfigSvc service.SystemConfigService,
	taskSvc service.PodcastTaskService,
) *ExecuteTaskCommand {
	return &ExecuteTaskCommand{
		ctx:             ctx,
		task:            task,
		systemConfigSvc: systemConfigSvc,
		taskSvc:         taskSvc,
	}
}

func (c *ExecuteTaskCommand) Execute(ctx context.Context) error {
	// config
	textAi, _, prompt, err := c.systemConfigSvc.GetPodcastConfig(ctx)
	if err != nil {
		return err
	}

	// execute task
	err = c.executeTaskState(ctx, textAi, prompt)
	if err != nil {
		logx.WithContext(c.ctx).Error("ExecuteTaskCommand.executeTaskState", err)
	}

	if err := c.taskSvc.SaveTask(c.ctx, c.task); err != nil {
		logx.WithContext(c.ctx).Error("ExecuteTaskCommand.SaveTask", err)
	}

	return err
}

// execute task state
func (c *ExecuteTaskCommand) executeTaskState(ctx context.Context, textAi *openai.Config, prompt *valueobject.PodcastScriptPrompt) error {
	var (
		messages = []*schema.Message{
			schema.SystemMessage(prompt.BuildSystemPrompt(c.task.Language)),
		}
		stylePrompt *valueobject.StylePrompt
	)

	for _, stage := range c.task.Stages {
		if !stage.Status.IsProcessing() {
			continue
		}

		userMsg := stage.BuildPrompt()

		if stage.Stage == valueobject.TaskStageStylize && stage.Prompt == "" {
			if stylePrompt == nil {
				stage.Fail("failed to classify news")
				c.task.Result = valueobject.TaskResultFailed

				return nil
			}

			stage.Prompt = stylePrompt.Prompt
			userMsg = stage.BuildPrompt()
		}

		messages = append(messages, schema.UserMessage(userMsg))

		stream, err := openai.NewChatModel(ctx, textAi).Stream(ctx, messages)
		if err != nil {
			stage.Fail(err.Error())
			c.task.Result = valueobject.TaskResultFailed
			return errors.WithStack(err)
		}
		defer stream.Close()

		var fullContent strings.Builder
		for {
			chunk, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				stage.Fail(err.Error())
				c.task.Result = valueobject.TaskResultFailed
				return errors.WithStack(err)
			}
			fullContent.WriteString(chunk.Content)
		}

		stage.SetOutput(fullContent.String())
		messages = append(messages, schema.AssistantMessage(stage.Output, nil))

		switch stage.Stage {
		case valueobject.TaskStageApproval:
			stage.SetStatus(prompt.VerifyApprovalResult(stage.Output))
		case valueobject.TaskStageClassify:
			stylePrompt = prompt.GetStylePrompt(stage.Output)
			if stylePrompt == nil {
				stage.Fail("failed to classify news")
			} else {
				stage.SetStatus(valueobject.StageStatusCompleted)
			}
		case valueobject.TaskStageScripted:
			logx.WithContext(ctx).Info("TaskStageScripted", stage.Output)
			scripts := prompt.ExtractScripts(stage.Output)
			if scripts == nil {
				stage.Fail("failed to extract scripts")
				stage.SetOutput(stage.Output)
			} else {
				stage.Audio.Scripts = scripts
				stage.SetStatus(valueobject.StageStatusCompleted)
			}
		default:
			stage.SetStatus(valueobject.StageStatusCompleted)
		}

		if stage.Status == valueobject.StageStatusFailed {
			c.task.Result = valueobject.TaskResultFailed

			return nil
		}
	}

	return nil
}
