package command

import (
	"context"

	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/service"
)

// RestyleArticleCommand restyle article command
type RestyleArticleCommand struct {
	ctx context.Context

	stageId uint
	prompt  string

	systemConfigSvc service.SystemConfigService
	taskSvc         service.PodcastTaskService
}

func NewRestyleArticleCommand(
	ctx context.Context,
	stageId uint,
	prompt string,
	systemConfigSvc service.SystemConfigService,
	taskSvc service.PodcastTaskService,
) *RestyleArticleCommand {
	return &RestyleArticleCommand{
		ctx:             ctx,
		stageId:         stageId,
		prompt:          prompt,
		systemConfigSvc: systemConfigSvc,
		taskSvc:         taskSvc,
	}
}

func (c *RestyleArticleCommand) Execute(ctx context.Context) error {
	// get config
	textAi, _, _, err := c.systemConfigSvc.GetPodcastConfig(ctx)
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

	// new stage
	var (
		stage = valueobject.NewTaskStage(valueobject.TaskStageStylize, c.prompt,
			valueobject.NewTaskAiFromTextAi(textAi))
		oldStage = task.GetStageById(c.stageId)
	)

	if oldStage.Stage == valueobject.TaskStageMerge {
		stage.Input = oldStage.Input
	} else {
		stage.Input = task.News.BuildPrompt()
	}

	task.AddNewStage(stage)

	if err = c.taskSvc.SaveTask(ctx, task); err != nil {
		return err
	}

	// execute task
	executeCmd := NewExecuteTaskCommand(c.ctx, task, c.systemConfigSvc, c.taskSvc)
	go func() {
		if err := executeCmd.Execute(executeCmd.ctx); err != nil {
			logx.WithContext(executeCmd.ctx).Error("RestyleArticleCommand", err)
		}
	}()

	return err
}
