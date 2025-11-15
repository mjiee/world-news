package command

import (
	"context"
	"strings"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/openai"
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
	textAi, _, prompt, err := getPodcastConfig(ctx, c.systemConfigSvc)
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

	// execute
	go c.restyleArticle(task, textAi, prompt)

	return err
}

// restyleArticle restyle article
func (c *RestyleArticleCommand) restyleArticle(task *entity.PodcastTask, textAi *openai.Config,
	prompt *valueobject.PodcastScriptPrompt) {
	var (
		stage    = task.GetStage(valueobject.TaskStageStylize)
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
		stage.SetStatus(valueobject.StageStatusCompleted)
	}

	if err := c.taskSvc.SaveTask(c.ctx, task); err != nil {
		logx.WithContext(c.ctx).Error("RestyleArticleCommand.SaveTask", err)
	}
}
