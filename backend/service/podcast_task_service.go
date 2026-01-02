package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/pkg/errors"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/audio"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/pathx"
	"github.com/mjiee/world-news/backend/pkg/ttsai"
	"github.com/mjiee/world-news/backend/repository"
	"github.com/mjiee/world-news/backend/repository/model"
)

// PodcastTaskService podcast task service
type PodcastTaskService interface {
	SaveTask(ctx context.Context, task *entity.PodcastTask) error
	DeleteTask(ctx context.Context, batchNo string) error
	EditScript(ctx context.Context, stageId uint, scripts []*ttsai.TtsScript) error
	UpdateTaskOutput(ctx context.Context, stageId uint, output string) error
	GetTaskByBatchNo(ctx context.Context, batchNo string) (*entity.PodcastTask, error)
	GetTaskByStageId(ctx context.Context, stageId uint) (*entity.PodcastTask, error)
	QueryTasks(ctx context.Context, params *valueobject.QueryPodcastTaskParams) ([]*entity.PodcastTask, int64, error)
	HasProcessingTasks(ctx context.Context, newsId uint) (bool, error)
	NewsHasTask(ctx context.Context, newsId uint) (bool, error)
	DownloadAudio(ctx context.Context, stageId uint, fileName string) error
}

type podcastTaskService struct {
}

func NewPodcastTaskService() PodcastTaskService {
	return &podcastTaskService{}
}

// SaveTask save task
func (s *podcastTaskService) SaveTask(ctx context.Context, task *entity.PodcastTask) error {
	data, err := task.ToModel()
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	if err := repository.Q.PodcastTask.WithContext(ctx).Save(data...); err != nil {
		return errors.WithStack(err)
	}

	task.Stages, err = gokit.SliceMapErr(data, valueobject.NewTaskStageFromModel)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTask delete task
func (s *podcastTaskService) DeleteTask(ctx context.Context, batchNo string) error {
	_, err := repository.Q.PodcastTask.WithContext(ctx).Where(repository.Q.PodcastTask.BatchNo.Eq(batchNo)).Delete()

	return errors.WithStack(err)
}

// EditScript edit script
func (s *podcastTaskService) EditScript(ctx context.Context, stageId uint, scripts []*ttsai.TtsScript) error {
	task, err := s.GetTaskByStageId(ctx, stageId)
	if err != nil {
		return err
	}

	stage := task.GetStageById(stageId)
	stage.Audio.Scripts = scripts

	return s.SaveTask(ctx, task)
}

// UpdateTaskOutput update task output
func (s *podcastTaskService) UpdateTaskOutput(ctx context.Context, stageId uint, output string) error {
	task, err := s.GetTaskByStageId(ctx, stageId)
	if err != nil {
		return err
	}

	stage := task.GetStageById(stageId)
	stage.SetOutput(output)

	return s.SaveTask(ctx, task)
}

// HasProcessingTasks check whether there are processing tasks
func (s *podcastTaskService) HasProcessingTasks(ctx context.Context, newsId uint) (bool, error) {
	repo := repository.Q.PodcastTask

	count, err := repo.WithContext(ctx).Where(repo.NewsId.Eq(newsId), repo.Result.Eq("")).Count()

	return count > 0, errors.WithStack(err)
}

// NewsHasTask check whether the news has task
func (s *podcastTaskService) NewsHasTask(ctx context.Context, newsId uint) (bool, error) {
	repo := repository.Q.PodcastTask

	count, err := repo.WithContext(ctx).Where(repo.NewsId.Eq(newsId)).Count()

	return count > 0, errors.WithStack(err)
}

// GetTaskByStageId get task by stage id
func (s *podcastTaskService) GetTaskByStageId(ctx context.Context, stageId uint) (*entity.PodcastTask, error) {
	repo := repository.Q.PodcastTask

	stage, err := repo.WithContext(ctx).Where(repo.ID.Eq(stageId)).First()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return s.GetTaskByBatchNo(ctx, stage.BatchNo)
}

// GetTaskByBatchNo get task by batch no
func (s *podcastTaskService) GetTaskByBatchNo(ctx context.Context, batchNo string) (*entity.PodcastTask, error) {
	repo := repository.Q.PodcastTask

	data, err := repo.WithContext(ctx).Where(repo.BatchNo.Eq(batchNo)).Find()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(data) == 0 {
		return nil, errorx.PodcastTaskNotFound
	}

	var (
		news     *model.NewsDetail
		newsRepo = repository.Q.NewsDetail
	)

	if data[0].NewsId > 0 {
		news, err = newsRepo.WithContext(ctx).Where(newsRepo.ID.Eq(data[0].NewsId)).First()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return entity.NewPodcastTaskFromModel(news, data)
}

// QueryTasks query tasks
func (s *podcastTaskService) QueryTasks(ctx context.Context, params *valueobject.QueryPodcastTaskParams) (
	[]*entity.PodcastTask, int64, error) {
	var (
		batchNos []string
		repo     = repository.Q.PodcastTask
	)

	total, err := repo.WithContext(ctx).Distinct(repo.BatchNo).Count()
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	if err = repo.WithContext(ctx).Group(repo.BatchNo).Order(repo.BatchNo.Desc()).
		Offset(params.Page.GetOffset()).Limit(params.Page.GetLimit()).
		Pluck(repo.BatchNo, &batchNos); err != nil {
		return nil, 0, errors.WithStack(err)
	}

	if len(batchNos) == 0 {
		return nil, total, nil
	}

	tasks, err := repo.WithContext(ctx).Where(repo.BatchNo.In(batchNos...)).Find()
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	var (
		newsRepo = repository.Q.NewsDetail
		newsIds  = gokit.SliceFilterMap(tasks, func(s *model.PodcastTask) (bool, uint) {
			return s.NewsId > 0, s.NewsId
		})
	)

	news, err := newsRepo.WithContext(ctx).Where(newsRepo.ID.In(newsIds...)).Find()
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	var (
		newsMap     = gokit.SliceToMap(news, func(s *model.NewsDetail) uint { return s.ID })
		tasksMap    = gokit.SliceGroupBy(tasks, func(t *model.PodcastTask) string { return t.BatchNo })
		taskNewsMap = make(map[string]*model.NewsDetail)
	)

	for _, task := range tasks {
		taskNewsMap[task.BatchNo] = newsMap[task.NewsId]
	}

	result := gokit.MapToSlice(tasksMap, func(batchNo string, tasks []*model.PodcastTask) *entity.PodcastTask {
		data, err := entity.NewPodcastTaskFromModel(taskNewsMap[batchNo], tasks)
		if err != nil {
			logx.WithContext(ctx).Error("QueryTasks.NewPodcastTaskFromModel", err)
			return nil
		}

		return data
	})

	result = gokit.SliceFilter(result, func(t *entity.PodcastTask) bool { return t != nil })

	slices.SortFunc(result, func(a, b *entity.PodcastTask) int {
		if a.CreatedAt.After(b.CreatedAt) {
			return -1
		} else {
			return 1
		}
	})

	return result, total, nil
}

// DownloadAudio download audio
func (s *podcastTaskService) DownloadAudio(ctx context.Context, stageId uint, fileName string) error {
	task, err := s.GetTaskByStageId(ctx, stageId)
	if err != nil {
		return err
	}

	stage := task.GetStageById(stageId)

	if stage.Audio == nil || stage.Audio.Data == "" {
		return errorx.PodcastTaskNotFound
	}

	if fileName == "" {
		fileName = fmt.Sprintf("%s_%d", task.BatchNo, stage.Id)
	}

	file := filepath.Join(pathx.GetDownloadPath(), fmt.Sprintf("%s.%s", fileName, stage.Audio.Type))

	data, err := base64.StdEncoding.DecodeString(stage.Audio.Data)
	if err != nil {
		return err
	}

	_, err = audio.WriteMp3sToFile(file, data)

	return errors.WithStack(err)
}
