package service

import (
	"context"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/httpx"
	"github.com/mjiee/world-news/backend/repository"
	"github.com/mjiee/world-news/backend/repository/model"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type NewsDetailService interface {
	CreateNews(ctx context.Context, news ...*entity.NewsDetail) error
	QueryNews(ctx context.Context, recordId uint, page *httpx.Pagination) ([]*entity.NewsDetail, int64, error)
	GetNewsDetail(ctx context.Context, id uint) (*entity.NewsDetail, error)
	DeleteNews(ctx context.Context, id uint) error
}

type newsDetailService struct {
}

func NewNewsDetailService() NewsDetailService {
	return &newsDetailService{}
}

// CreateNews creates a new news detail.
func (s *newsDetailService) CreateNews(ctx context.Context, news ...*entity.NewsDetail) error {
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
func (s *newsDetailService) QueryNews(ctx context.Context, recordId uint, page *httpx.Pagination) ([]*entity.NewsDetail, int64, error) {
	repo := repository.Q.NewsDetail

	data, total, err := repo.WithContext(ctx).Where(repo.RecordId.Eq(recordId)).FindByPage(page.GetOffset(), page.GetLimit())
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
func (s *newsDetailService) GetNewsDetail(ctx context.Context, id uint) (*entity.NewsDetail, error) {
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
func (s *newsDetailService) DeleteNews(ctx context.Context, id uint) error {
	_, err := repository.Q.NewsDetail.WithContext(ctx).Where(repository.Q.NewsDetail.ID.Eq(id)).Delete()

	return errors.WithStack(err)
}
