package service

import (
	"context"
	"errors"

	"github.com/gocolly/colly/v2"
	"gorm.io/gorm"

	"github.com/mjiee/world-news/backend/repository"
)

type NewsCrawlingService interface {
}

type newsCrawlingService struct {
	collector          *colly.Collector
	newsDetailRepo     repository.NewsDetailRepository
	crawlingRecordRepo repository.CrawlingRecordRepository
}

func NewNewsCrawlingService(c *colly.Collector, db *gorm.DB) NewsCrawlingService {
	return &newsCrawlingService{collector: c,
		newsDetailRepo:     repository.NewNewsDetailRepository(db),
		crawlingRecordRepo: repository.NewCrawlingRecordRepository(db),
	}
}

// getCollector 获取Collector
func (s *newsCrawlingService) getCollector() *colly.Collector {
	return s.collector.Clone()
}

// collectorIgnoreErr 可忽略错误
func collectorIgnoreErr(err error) bool {
	return errors.Is(err, colly.ErrAlreadyVisited) || errors.Is(err, context.DeadlineExceeded)
}
