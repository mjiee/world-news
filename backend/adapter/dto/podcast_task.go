package dto

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/httpx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/ttsai"
)

// CreateTaskRequest is the request for creating a podcast task
type CreateTaskRequest struct {
	Language string      `json:"language" binding:"required"`
	News     *NewsDetail `json:"news" binding:"required"`
	VoiceIds []string    `json:"voiceIds,omitempty"`
}

// CreateTaskResult is the result for creating a podcast task
type CreateTaskResult struct {
	BatchNo string `json:"batchNo"`
}

// DeleteTaskRequest is the request for deleting a podcast task
type DeleteTaskRequest struct {
	BatchNo string `json:"batchNo" binding:"required"`
}

// RestyleArticleRequest is the request for restyling an article
type RestyleArticleRequest struct {
	StageId uint   `json:"stageId" binding:"required"`
	Prompt  string `json:"prompt" binding:"required"`
}

// CreateScriptRequest is the request for creating a podcast script
type CreateScriptRequest struct {
	StageId  uint     `json:"stageId" binding:"required"`
	VoiceIds []string `json:"voiceIds"`
}

// EditScript is the request for editing a podcast script
type EditScriptRequest struct {
	StageId uint               `json:"stageId" binding:"required"`
	Scripts []*ttsai.TtsScript `json:"scripts"`
}

// CreateAudioRequest is the request for generating a podcast audio
type CreateAudioRequest struct {
	StageId uint `json:"stageId" binding:"required"`
}

// DownloadAudioRequest is the request for downloading a podcast audio
type DownloadAudioRequest struct {
	StageId uint `json:"stageId" binding:"required"`
}

// GetTaskRequest is the request for getting a podcast task
type GetTaskRequest struct {
	BatchNo string `json:"batchNo" binding:"required"`
}

// QueryTaskRequest is the request for querying podcast tasks
type QueryTaskRequest struct {
	StartDate  string            `json:"startDate,omitempty"`
	EndDate    string            `json:"endDate,omitempty"`
	Pagination *httpx.Pagination `json:"pagination"`
}

// ToValueObject converts the request to a value object
func (r *QueryTaskRequest) ToValueObject() *valueobject.QueryPodcastTaskParams {
	query := &valueobject.QueryPodcastTaskParams{Page: r.Pagination}

	if r.StartDate != "" {
		createDate, err := time.Parse(time.DateOnly, r.StartDate)
		if err != nil {
			logx.Error("parse start date error", errors.WithStack(err))
		} else {
			query.StartDate = createDate
		}
	}

	if r.EndDate != "" {
		createDate, err := time.Parse(time.DateOnly, r.EndDate)
		if err != nil {
			logx.Error("parse end date error", errors.WithStack(err))
		} else {
			query.EndDate = createDate
		}
	}

	return query
}

// QueryTaskResult is the result for querying podcast tasks
type QueryTaskResult struct {
	Data  []*PodcastTask `json:"data"`
	Total int64          `json:"total"`
}

func NewQueryTaskResult(tasks []*entity.PodcastTask, total int64) *QueryTaskResult {
	return &QueryTaskResult{
		Data: gokit.SliceMap(tasks, func(t *entity.PodcastTask) *PodcastTask {
			task := NewPodcastTask(t)
			task.Stages = nil

			return task
		}),
		Total: total,
	}
}

// PodcastTask is the podcast task
type PodcastTask struct {
	BatchNo   string       `json:"batchNo"`
	Title     string       `json:"title"`
	News      *NewsDetail  `json:"news,omitempty"`
	Language  string       `json:"language"`
	Result    string       `json:"result"`
	Stages    []*TaskStage `json:"stages,omitempty"`
	CreatedAt string       `json:"createdAt"`
}

// NewPodcastTask converts the entity to a podcast task
func NewPodcastTask(task *entity.PodcastTask) *PodcastTask {
	data := &PodcastTask{
		BatchNo:   task.BatchNo,
		Title:     task.Title,
		Language:  task.Language,
		Result:    string(task.Result),
		Stages:    gokit.SliceMap(task.Stages, NewTaskStage),
		CreatedAt: task.CreatedAt.Format(time.DateTime),
	}

	if task.News != nil {
		data.News = NewBaseNewsDetail(task.News)
		if data.Title == "" {
			data.Title = task.News.Title
		}
	}

	return data
}

// TaskStage is the task stage
type TaskStage struct {
	Id        uint                      `json:"id"`
	BatchNo   string                    `json:"batchNo"`
	Stage     string                    `json:"stage"`
	Status    string                    `json:"status"`
	Prompt    string                    `json:"prompt"`
	Output    string                    `json:"output"`
	Reason    string                    `json:"reason"`
	Audio     *valueobject.PodcastAudio `json:"audio"`
	TaskAi    *valueobject.TaskAi       `json:"taskAi"`
	CreatedAt string                    `json:"createdAt"`
	UpdatedAt string                    `json:"updatedAt"`
}

// NewTaskStage converts the entity to a task stage
func NewTaskStage(stage *valueobject.TaskStage) *TaskStage {
	if stage == nil {
		return nil
	}

	task := &TaskStage{
		Id:        stage.Id,
		BatchNo:   stage.BatchNo,
		Stage:     string(stage.Stage),
		Status:    string(stage.Status),
		Prompt:    stage.Prompt,
		Output:    stage.Output,
		Reason:    stage.Reason,
		CreatedAt: stage.CreatedAt.Format(time.DateTime),
		UpdatedAt: stage.UpdatedAt.Format(time.DateTime),
	}

	if stage.Audio != nil {
		task.Audio = stage.Audio
	}

	if stage.TaskAi != nil {
		task.TaskAi = stage.TaskAi
	}

	return task
}

// MergeArticleRequest is the request for merging an article
type MergeArticleRequest struct {
	Language string   `json:"language" binding:"required"`
	Title    string   `json:"title" binding:"required"`
	StageIds []uint   `json:"stageIds" binding:"required"`
	VoiceIds []string `json:"voiceIds,omitempty"`
}
