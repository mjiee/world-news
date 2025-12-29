package command

import (
	"context"
	"errors"
	"slices"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/openai"
	"github.com/mjiee/world-news/backend/pkg/ttsai"
	"github.com/mjiee/world-news/backend/service"
)

// CreateTaskCommand represents a command to create a news to podcast task.
type CreateTaskCommand struct {
	ctx      context.Context
	language string
	news     *entity.NewsDetail
	voiceIds []string

	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
	taskSvc         service.PodcastTaskService
}

func NewCreateTaskCommand(
	ctx context.Context,
	language string,
	news *entity.NewsDetail,
	voiceIds []string,
	newsSvc service.NewsService,
	systemConfigSvc service.SystemConfigService,
	taskSvc service.PodcastTaskService,
) *CreateTaskCommand {
	return &CreateTaskCommand{
		ctx:             ctx,
		language:        language,
		news:            news,
		voiceIds:        voiceIds,
		newsSvc:         newsSvc,
		systemConfigSvc: systemConfigSvc,
		taskSvc:         taskSvc,
	}
}

func (c *CreateTaskCommand) Execute(ctx context.Context) (string, error) {
	if c.news == nil {
		return "", errorx.ParamsError
	}

	// config
	textAi, ttsAi, prompt, err := c.systemConfigSvc.GetPodcastConfig(ctx)
	if err != nil {
		return "", err
	}

	newTask, err := c.createTask(ctx, textAi, ttsAi, prompt)
	if err != nil {
		return "", err
	}

	// execute task
	executeCmd := NewExecuteTaskCommand(c.ctx, newTask, c.systemConfigSvc, c.taskSvc)
	go func() {
		if err := executeCmd.Execute(executeCmd.ctx); err != nil {
			logx.WithContext(executeCmd.ctx).Error("CreateTaskCommand", err)
		}
	}()

	return newTask.BatchNo, nil
}

// create new task
func (c *CreateTaskCommand) createTask(ctx context.Context, textAi *openai.Config, ttsAi *ttsai.Config,
	prompt *valueobject.PodcastScriptPrompt) (*entity.PodcastTask, error) {
	voices := gokit.SliceFilter(ttsAi.Voices, func(v *ttsai.Voice) bool { return slices.Contains(c.voiceIds, v.Id) })
	if len(ttsAi.Voices) == 0 || len(voices) != len(c.voiceIds) {
		return nil, errorx.PodcastVoiceNotFound
	}

	if c.news == nil {
		return nil, errorx.ParamsError
	}

	// news
	news, err := c.newsSvc.GetNewsDetail(ctx, c.news.Id)
	if err != nil && !errors.Is(err, errorx.NewsNotFound) {
		return nil, err
	}

	if news == nil {
		if err := c.newsSvc.CreateNews(ctx, c.news); err != nil {
			return nil, err
		}
	}

	// check crawling record
	hasProcessingTask, err := c.taskSvc.HasProcessingTasks(ctx, c.news.Id)
	if err != nil {
		return nil, err
	}

	if hasProcessingTask {
		return nil, errorx.HasProcessingTasks
	}

	// create task
	newTask := entity.NewPodcastTask(c.news, c.language)

	if err := c.buildTaskState(newTask, textAi, ttsAi, prompt); err != nil {
		return nil, err
	}

	err = c.taskSvc.SaveTask(c.ctx, newTask)

	return newTask, err
}

// build task state
func (c *CreateTaskCommand) buildTaskState(task *entity.PodcastTask, textAi *openai.Config, ttsAi *ttsai.Config,
	prompt *valueobject.PodcastScriptPrompt) error {
	ai := valueobject.NewTaskAiFromTextAi(textAi)

	for _, stage := range valueobject.StagePriority {
		switch stage {
		case valueobject.TaskStageApproval:
			if prompt.ApprovalPrompt == "" {
				continue
			}

			stage := valueobject.NewTaskStage(stage, prompt.BuildApprovalPrompt(task.Language), ai)
			stage.Input = task.News.BuildPrompt()

			task.AddNewStage(stage)
		case valueobject.TaskStageRewrite:
			if prompt.RewritePrompt == "" {
				continue
			}

			stage := valueobject.NewTaskStage(stage, prompt.RewritePrompt, ai)
			if len(task.Stages) == 0 {
				stage.Input = task.News.BuildPrompt()
			}

			task.AddNewStage(stage)
		case valueobject.TaskStageClassify:
			if len(prompt.StylizePrompts) == 0 {
				continue
			}

			// classify
			stage := valueobject.NewTaskStage(stage, prompt.BuildClassifyPrompt(task.Language), ai)
			if len(task.Stages) == 0 {
				stage.Input = task.News.BuildPrompt()
			}

			task.AddNewStage(stage)

			// stylize
			stylizeStage := valueobject.NewTaskStage(valueobject.TaskStageStylize, "", ai)
			task.AddNewStage(stylizeStage)
		case valueobject.TaskStageScripted:
			if len(c.voiceIds) == 0 {
				if len(task.Stages) != 0 {
					continue
				}

				c.voiceIds = []string{ttsAi.Voices[0].Id}
			}

			voices := gokit.SliceFilter(ttsAi.Voices, func(v *ttsai.Voice) bool {
				return slices.Contains(c.voiceIds, v.Id)
			})

			stage := valueobject.NewTaskStage(stage, valueobject.BuildScriptPrompt(task.Language, voices), ai)

			if len(task.Stages) == 0 {
				stage.Input = task.News.BuildPrompt()
			}

			stage.Audio = &valueobject.PodcastAudio{Voices: voices}
			task.AddNewStage(stage)
		}
	}

	return nil
}
