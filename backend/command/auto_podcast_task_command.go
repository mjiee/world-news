package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/ttsai"
	"github.com/mjiee/world-news/backend/service"
)

// AutoPodcastTaskCommand represents the command to auto podcast task.
type AutoPodcastTaskCommand struct {
	ctx      context.Context
	language string
	news     *entity.NewsDetail

	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
	taskSvc         service.PodcastTaskService
}

func NewAutoPodcastTaskCommand(
	ctx context.Context,
	language string,
	news *entity.NewsDetail,
	newsSvc service.NewsService,
	systemConfigSvc service.SystemConfigService,
	taskSvc service.PodcastTaskService,
) *AutoPodcastTaskCommand {
	return &AutoPodcastTaskCommand{
		ctx:             ctx,
		language:        language,
		news:            news,
		newsSvc:         newsSvc,
		systemConfigSvc: systemConfigSvc,
		taskSvc:         taskSvc,
	}
}

func (c *AutoPodcastTaskCommand) Execute(ctx context.Context) (string, error) {
	// 获取系统配置
	textAiConfig, ttsAiConfig, prompt, err := c.systemConfigSvc.GetPodcastConfig(c.ctx)
	if err != nil {
		return "", err
	}

	if len(ttsAiConfig.Voices) == 0 {
		return "", nil
	}

	voiceIds := gokit.SliceMap(ttsAiConfig.Voices, func(v *ttsai.Voice) string { return v.Id })
	if len(voiceIds) > 2 {
		voiceIds = voiceIds[:2]
	}

	// 创建任务
	createCmd := NewCreateTaskCommand(c.ctx, c.language, c.news, voiceIds, c.newsSvc, c.systemConfigSvc, c.taskSvc)

	newTask, err := createCmd.createTask(ctx, textAiConfig, ttsAiConfig, prompt)
	if err != nil {
		return "", err
	}

	go c.executeTask(newTask)

	return newTask.BatchNo, err
}

func (c *AutoPodcastTaskCommand) executeTask(task *entity.PodcastTask) {
	// execute task
	executeCmd := NewExecuteTaskCommand(c.ctx, task, c.systemConfigSvc, c.taskSvc)
	if err := executeCmd.Execute(c.ctx); err != nil {
		logx.Error("AutoPodcastTaskCommand.executeTask", err)
		return
	}

	scriptedStage := task.GetStage(valueobject.TaskStageScripted)
	if scriptedStage == nil {
		logx.Error("AutoPodcastTaskCommand.executeTask", errors.New("scriptedStage is nil"))
		return
	}

	// create audio
	createAudioCmd := NewCreateAudioCommand(c.ctx, scriptedStage.Id, c.systemConfigSvc, c.taskSvc)
	if err := createAudioCmd.Execute(c.ctx); err != nil {
		logx.Error("AutoPodcastTaskCommand.executeTask", err)
	}
}
