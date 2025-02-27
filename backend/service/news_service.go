package service

import (
	"context"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/repository"
	"github.com/mjiee/world-news/backend/repository/model"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// NewsService represents the interface for news-related operations.
type NewsService interface {
	CreateNews(ctx context.Context, news ...*entity.NewsDetail) error
	QueryNews(ctx context.Context, params *valueobject.QueryNewsParams) ([]*entity.NewsDetail, int64, error)
	GetNewsDetail(ctx context.Context, id uint) (*entity.NewsDetail, error)
	DeleteNews(ctx context.Context, id uint) error
}

type newsService struct {
}

func NewNewsService() NewsService {
	return &newsService{}
}

// CreateNews creates a new news detail.
func (s *newsService) CreateNews(ctx context.Context, news ...*entity.NewsDetail) error {
	if len(news) == 0 {
		return nil
	}

	data := make([]*model.NewsDetail, len(news))

	for i, v := range news {
		m, err := v.ToModel()
		if err != nil {
			return errors.WithStack(err)
		}

		data[i] = m
	}

	return errors.WithStack(repository.Q.NewsDetail.WithContext(ctx).CreateInBatches(data, 100))
}

// QueryNews queries news details based on the provided record ID and pagination.
func (s *newsService) QueryNews(ctx context.Context, params *valueobject.QueryNewsParams) (
	[]*entity.NewsDetail, int64, error) {
	var (
		repo  = repository.Q.NewsDetail
		query = repo.WithContext(ctx)
	)

	if params.RecordId != 0 {
		query = query.Where(repo.RecordId.Eq(params.RecordId))
	}

	if params.Source != "" {
		query = query.Where(repo.Source.Eq(params.Source))
	}

	if params.Topic != "" {
		query = query.Where(repo.Topic.Eq(params.Topic))
	}

	data, total, err := query.Order(repo.ID.Desc()).FindByPage(params.Page.GetOffset(), params.Page.GetLimit())
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	news := make([]*entity.NewsDetail, len(data))

	for i, v := range data {
		news[i], err = entity.NewNewsDetailFromModel(v)
		if err != nil {
			return nil, 0, errors.WithStack(err)
		}
	}

	return news, total, nil
}

// GetNewsDetail retrieves the news detail based on the provided ID.
func (s *newsService) GetNewsDetail(ctx context.Context, id uint) (*entity.NewsDetail, error) {
	repo := repository.Q.NewsDetail

	news, err := repo.WithContext(ctx).Where(repo.ID.Eq(id)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.NewsNotFound
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.NewNewsDetailFromModel(news)
}

// DeleteNews deletes the news detail based on the provided ID.
func (s *newsService) DeleteNews(ctx context.Context, id uint) error {
	_, err := repository.Q.NewsDetail.WithContext(ctx).Where(repository.Q.NewsDetail.ID.Eq(id)).Delete()

	return errors.WithStack(err)
}
