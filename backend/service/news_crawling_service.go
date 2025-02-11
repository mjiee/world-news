package service

import (
	"context"
	"errors"

	"github.com/gocolly/colly/v2"
)

type NewsCrawlingService interface {
}

type newsCrawlingService struct {
	collector *colly.Collector
}

func NewNewsCrawlingService(c *colly.Collector) NewsCrawlingService {
	return &newsCrawlingService{collector: c}
}

// getCollector 获取Collector
func (s *newsCrawlingService) getCollector() *colly.Collector {
	return s.collector.Clone()
}

// collectorIgnoreErr 可忽略错误
func collectorIgnoreErr(err error) bool {
	return errors.Is(err, colly.ErrAlreadyVisited) || errors.Is(err, context.DeadlineExceeded)
}
