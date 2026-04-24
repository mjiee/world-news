package entity

import (
	"slices"
	"time"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/ttsai"
	"github.com/mjiee/world-news/backend/repository/model"
)

// PodcastTask represents the detailed information about a news to podcast task.
type PodcastTask struct {
	BatchNo   string
	Title     string
	News      *NewsDetail
	Language  string
	Stages    []*valueobject.TaskStage
	Result    valueobject.TaskResult
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewPodcastTask creates a new PodcastTask instance.
func NewPodcastTask(news *NewsDetail, language string) *PodcastTask {
	return &PodcastTask{
		News:     news,
		Language: language,
		BatchNo:  time.Now().Format("060102150405"),
		Stages:   make([]*valueobject.TaskStage, 0),
	}
}

// NewPodcastTaskFromModel converts a model.PodcastTask to a PodcastTask.
func NewPodcastTaskFromModel(news *model.NewsDetail, tasks []*model.PodcastTask) (*PodcastTask, error) {
	slices.SortFunc(tasks, func(a, b *model.PodcastTask) int { return int(a.ID - b.ID) })

	data := &PodcastTask{}

	for _, task := range tasks {
		data.BatchNo = task.BatchNo

		stage, err := valueobject.NewTaskStageFromModel(task)
		if err != nil {
			return nil, err
		}

		data.Stages = append(data.Stages, stage)

		if data.CreatedAt.IsZero() || data.CreatedAt.After(task.CreatedAt) {
			data.CreatedAt = task.CreatedAt
		}

		if data.CreatedAt.IsZero() || data.UpdatedAt.Before(task.UpdatedAt) {
			data.UpdatedAt = task.UpdatedAt
		}

		if task.Result != "" {
			data.Result = valueobject.TaskResult(task.Result)
		}

		if task.Language != "" {
			data.Language = task.Language
		}

		if task.Title != "" {
			data.Title = task.Title
		} else if data.Title == "" && news != nil {
			data.Title = news.Title
		}
	}

	if news != nil {
		newsDetail, err := NewNewsDetailFromModel(news)
		if err != nil {
			return nil, err
		}

		data.News = newsDetail
	}

	return data, nil
}

// ToModel converts the PodcastTask to a model.PodcastTask.
func (t *PodcastTask) ToModel() ([]*model.PodcastTask, error) {
	return gokit.SliceMapErr(t.Stages, func(i *valueobject.TaskStage) (*model.PodcastTask, error) {
		if i.CreatedAt.IsZero() {
			i.CreatedAt = time.Now()
		}

		if i.UpdatedAt.IsZero() {
			i.UpdatedAt = time.Now()
		}

		data := &model.PodcastTask{
			ID:        i.Id,
			BatchNo:   t.BatchNo,
			Title:     t.Title,
			Stage:     string(i.Stage),
			Status:    string(i.Status),
			Language:  t.Language,
			Prompt:    i.Prompt,
			Input:     i.Input,
			Output:    i.Output,
			Reason:    i.Reason,
			Result:    string(t.Result),
			Audio:     gokit.MarshalSafe(i.Audio),
			TaskAi:    gokit.MarshalSafe(i.TaskAi),
			Extra:     gokit.MarshalSafe(i.Extra),
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
		}

		if t.News != nil {
			data.NewsId = t.News.Id
		}

		return data, nil
	})
}

// AddNewStage adds a new task stage to the PodcastTask.
func (t *PodcastTask) AddNewStage(stage *valueobject.TaskStage) {
	t.Stages = append(t.Stages, stage)
	t.Result = ""
}

// GetPodcastScript gets the podcast script for the PodcastTask.
func (t *PodcastTask) GetPodcastScript() []*ttsai.TtsScript {
	stage := gokit.SliceFindLast(t.Stages, func(t *valueobject.TaskStage) bool {
		return t.Stage == valueobject.TaskStageScripted
	})

	if stage.Audio == nil {
		return nil
	}

	return stage.Audio.Scripts
}

// GetStageById gets a task stage by its ID.
func (t *PodcastTask) GetStageById(id uint) *valueobject.TaskStage {
	return gokit.SliceFind(t.Stages, func(t *valueobject.TaskStage) bool {
		return t.Id == id
	})
}

// GetStage get a task stage.
func (t *PodcastTask) GetStage(stage valueobject.TaskStageName) *valueobject.TaskStage {
	return gokit.SliceFindLast(t.Stages, func(t *valueobject.TaskStage) bool {
		return t.Stage == stage
	})
}

// VerifyTask verify task
func (t *PodcastTask) VerifyTask() error {
	if len(t.Stages) == 0 {
		return errorx.PodcastTaskNotFound
	}

	stage := gokit.SliceFindLast(t.Stages, func(t *valueobject.TaskStage) bool {
		return t.Status == valueobject.StageStatusProcessing
	})
	if stage != nil {
		return errorx.HasProcessingTasks
	}

	return nil
}
