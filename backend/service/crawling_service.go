package service

import (
	"context"

	"github.com/mjiee/world-news/backend/entity"
	"github.com/mjiee/world-news/backend/pkg/httpx"
	"github.com/mjiee/world-news/backend/repository"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
)

// CrawlingService news crawling service
type CrawlingService interface {
	GetCollector() *colly.Collector
	CreateCrawlingRecord(ctx context.Context, record *entity.CrawlingRecord) error
	UpdateCrawlingRecord(ctx context.Context, record *entity.CrawlingRecord) error
	QueryCrawlingRecords(ctx context.Context, page *httpx.Pagination) ([]*entity.CrawlingRecord, int64, error)
	DeleteCrawlingRecord(ctx context.Context, id uint) error
}

type crawlingService struct {
	collector *colly.Collector
}

func NewCrawlingService(c *colly.Collector) CrawlingService {
	return &crawlingService{collector: c}
}

// GetCollector get a new collector
func (s *crawlingService) GetCollector() *colly.Collector {
	return s.collector.Clone()
}

// CreateCrawlingRecord create crawling record
func (s *crawlingService) CreateCrawlingRecord(ctx context.Context, record *entity.CrawlingRecord) error {
	data, err := record.ToModel()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := repository.Q.CrawlingRecord.WithContext(ctx).Create(data); err != nil {
		return errors.WithStack(err)
	}

	record.Id = data.Id

	return nil
}

// UpdateCrawlingRecord update crawling record status
func (s *crawlingService) UpdateCrawlingRecord(ctx context.Context, record *entity.CrawlingRecord) error {
	data, err := record.ToModel()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = repository.Q.CrawlingRecord.WithContext(ctx).Updates(data)

	return errors.WithStack(err)
}

// QueryCrawlingRecords get crawling records
func (s *crawlingService) QueryCrawlingRecords(ctx context.Context, page *httpx.Pagination) ([]*entity.CrawlingRecord, int64, error) {
	data, total, err := repository.Q.CrawlingRecord.WithContext(ctx).FindByPage(page.GetOffset(), page.GetLimit())
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	records := make([]*entity.CrawlingRecord, len(data))

	for idx, v := range data {
		records[idx], err = entity.NewCrawlingRecordFromModel(v)
		if err != nil {
			return nil, 0, errors.WithStack(err)
		}
	}

	return records, total, nil
}

// DeleteCrawlingRecord delete crawling record
func (s *crawlingService) DeleteCrawlingRecord(ctx context.Context, id uint) error {
	err := repository.Q.Transaction(func(tx *repository.Query) error {
		if _, err := tx.CrawlingRecord.WithContext(ctx).Where(tx.CrawlingRecord.Id.Eq(uint(id))).Delete(); err != nil {
			return errors.WithStack(err)
		}

		if _, err := tx.NewsDetail.WithContext(ctx).Where(tx.NewsDetail.RecordId.Eq(id)).Delete(); err != nil {
			return errors.WithStack(err)
		}

		return nil
	})

	return errors.WithStack(err)
}
