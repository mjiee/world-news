package command

import (
	"context"
	"fmt"
	"slices"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/ttsai"
	"github.com/mjiee/world-news/backend/service"
)

// MergeArticleCommand is the command to merge article
type MergeArticleCommand struct {
	ctx      context.Context
	language string
	title    string
	stageIds []uint
	voiceIds []string

	systemConfigSvc service.SystemConfigService
	taskSvc         service.PodcastTaskService
}

func NewMergeArticleCommand(
	ctx context.Context,
	language string,
	title string,
	stageIds []uint,
	voiceIds []string,
	systemConfigSvc service.SystemConfigService,
	taskSvc service.PodcastTaskService,
) *MergeArticleCommand {
	return &MergeArticleCommand{
		language:        language,
		title:           title,
		stageIds:        stageIds,
		voiceIds:        voiceIds,
		systemConfigSvc: systemConfigSvc,
		taskSvc:         taskSvc,
	}
}

func (c *MergeArticleCommand) Execute(ctx context.Context) (string, error) {
	if len(c.stageIds) == 0 || c.title == "" {
		return "", errorx.ParamsError
	}

	// config
	textAi, ttsAi, prompt, err := c.systemConfigSvc.GetPodcastConfig(ctx)
	if err != nil {
		return "", err
	}

	voices := gokit.SliceFilter(ttsAi.Voices, func(v *ttsai.Voice) bool { return slices.Contains(c.voiceIds, v.Id) })
	if len(voices) != len(c.voiceIds) {
		return "", errorx.PodcastVoiceNotFound
	}

	// get task
	tasks, err := gokit.SliceMapErr(c.stageIds, func(id uint) (*entity.PodcastTask, error) {
		return c.taskSvc.GetTaskByStageId(ctx, id)
	})
	if err != nil {
		return "", err
	}

	contents, err := gokit.SliceMapErr(c.stageIds, func(id uint) (string, error) {
		for _, task := range tasks {
			stage := task.GetStageById(id)
			if stage != nil {
				return stage.Output, nil
			}
		}

		return "", errorx.PodcastTaskNotFound
	})
	if err != nil {
		return "", err
	}

	// new task
	var (
		task = entity.NewPodcastTask(nil, c.language)
		ai   = valueobject.NewTaskAiFromTextAi(textAi)
	)

	task.Title = c.title

	// merge stage
	mergeStage := valueobject.NewTaskStage(valueobject.TaskStageMerge, prompt.BuildMergePrompt(c.language), ai)

	for idx, content := range contents {
		mergeStage.Input = fmt.Sprintf("%s\n\nThe %d'st podcast: \n%s", mergeStage.Input, idx+1, content)
	}

	mergeStage.Extra = &valueobject.TaskStageExtra{
		NewsIds: gokit.SliceFilterMap(tasks, func(t *entity.PodcastTask) (bool, uint) {
			if t.News == nil {
				return false, 0
			}

			return true, t.News.Id
		}),
		StageIds: c.stageIds,
	}

	task.AddNewStage(mergeStage)

	// scritp stage
	if len(c.voiceIds) > 0 {
		scriptStage := valueobject.NewTaskStage(valueobject.TaskStageScripted,
			valueobject.BuildScriptPrompt(task.Language, voices), ai)

		scriptStage.Audio = &valueobject.PodcastAudio{Voices: voices}
		task.AddNewStage(scriptStage)
	}

	if err := c.taskSvc.SaveTask(ctx, task); err != nil {
		return "", err
	}

	// execute task
	executeCmd := NewExecuteTaskCommand(c.ctx, task, c.systemConfigSvc, c.taskSvc)
	go func() {
		if err := executeCmd.Execute(executeCmd.ctx); err != nil {
			logx.WithContext(executeCmd.ctx).Error("MergeArticleCommand", err)
		}
	}()

	return task.BatchNo, nil
}
