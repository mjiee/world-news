package command

import (
	"context"
	"errors"
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
	textAi, ttsAi, prompt, err := getPodcastConfig(ctx, c.systemConfigSvc)
	if err != nil {
		return "", err
	}

	voices := gokit.SliceFilter(ttsAi.Voices, func(v *ttsai.Voice) bool { return slices.Contains(c.voiceIds, v.Id) })
	if len(ttsAi.Voices) == 0 || len(voices) != len(c.voiceIds) {
		return "", errorx.PodcastVoiceNotFound
	}

	// news
	news, err := c.newsSvc.GetNewsDetail(ctx, c.news.Id)
	if err != nil && !errors.Is(err, errorx.NewsNotFound) {
		return "", err
	}

	if news == nil {
		if err := c.newsSvc.CreateNews(ctx, c.news); err != nil {
			return "", err
		}
	}

	// check crawling record
	hasProcessingTask, err := c.taskSvc.HasProcessingTasks(ctx, c.news.Id)
	if err != nil {
		return "", err
	}

	if hasProcessingTask {
		return "", errorx.HasProcessingTasks
	}

	// create task
	newTask := entity.NewPodcastTask(c.news, c.language)

	if err := c.buildTaskState(newTask, textAi, ttsAi, prompt); err != nil {
		return "", err
	}

	if err := c.taskSvc.SaveTask(c.ctx, newTask); err != nil {
		return "", err
	}

	// execute task
	go c.executeTaskState(newTask, textAi, prompt)

	return newTask.BatchNo, nil
}

func getPodcastConfig(ctx context.Context, configSvc service.SystemConfigService) (
	textAi *openai.Config,
	ttsAi *ttsai.Config,
	prompt *valueobject.PodcastScriptPrompt,
	err error,
) {
	// text ai config
	config, err := configSvc.GetSystemConfig(ctx, valueobject.TextAIKey.String())
	if err != nil {
		return nil, nil, nil, err
	}

	textAi, err = entity.UnmarshalValue[openai.Config](config, errorx.OpenaiConfigNotFound)
	if err != nil {
		return nil, nil, nil, err
	}

	// tts ai config
	ttsAiConfig, err := configSvc.GetSystemConfig(ctx, valueobject.TextToSpeechAIKey.String())
	if err != nil {
		return nil, nil, nil, err
	}

	ttsAi, err = entity.UnmarshalValue[ttsai.Config](ttsAiConfig, errorx.TtsAiConfigNotFound)
	if err != nil {
		return nil, nil, nil, err
	}

	// podcast script prompt
	config, err = configSvc.GetSystemConfig(ctx, valueobject.PodcastScriptPromptKey.String())
	if err != nil {
		return nil, nil, nil, err
	}

	prompt, err = entity.UnmarshalValue[valueobject.PodcastScriptPrompt](config, errorx.PodcastPromptNotFound)
	if err != nil {
		return nil, nil, nil, err
	}

	return textAi, ttsAi, prompt, nil
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

// executeTaskState execute task state
func (c *CreateTaskCommand) executeTaskState(
	task *entity.PodcastTask,
	textAi *openai.Config,
	prompt *valueobject.PodcastScriptPrompt,
) {
	// execute task
	if err := executeTaskState(c.ctx, task, textAi, prompt); err != nil {
		logx.WithContext(c.ctx).Error("CreatePodcastTaskCommand.executeTaskState", err)
	}

	if err := c.taskSvc.SaveTask(c.ctx, task); err != nil {
		logx.WithContext(c.ctx).Error("CreatePodcastTaskCommand.SaveTask", err)
	}
}

// execute task state
func executeTaskState(ctx context.Context, newTask *entity.PodcastTask, textAi *openai.Config,
	prompt *valueobject.PodcastScriptPrompt) error {
	var (
		messages    = []*openai.Message{openai.SystemMessage(prompt.BuildSystemPrompt(newTask.Language))}
		stylePrompt *valueobject.StylePrompt
	)

	for _, stage := range newTask.Stages {
		userMsg := stage.BuildPrompt()

		if stage.Stage == valueobject.TaskStageStylize {
			if stylePrompt == nil {
				stage.Fail("failed to classify news")
				newTask.Result = valueobject.TaskResultFailed

				return nil
			}

			stage.Prompt = stylePrompt.Prompt
			userMsg = stage.BuildPrompt()
		}

		messages = append(messages, openai.UserMessage(userMsg))

		data, err := openai.NewOpenaiClient(textAi).SetMessage(messages...).ChatCompletion(ctx)
		if err != nil {
			stage.Fail(err.Error())
			newTask.Result = valueobject.TaskResultFailed

			return err
		}

		assistantMsg := gokit.SliceMap(data.Choices, func(item *openai.ChatCompletionChoice) string {
			return item.Message.Content
		})

		stage.TaskAi.SessionId = data.ID
		stage.SetOutput(strings.Join(assistantMsg, "\n"))
		messages = append(messages, openai.AssistantMessage(stage.Output))

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
			newTask.Result = valueobject.TaskResultFailed

			return nil
		}
	}

	return nil
}
