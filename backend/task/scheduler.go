//go:build web

package task

import (
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/service"

	"github.com/go-co-op/gocron/v2"
	"github.com/pkg/errors"
)

// scheduler is a struct that manages scheduled tasks for the application.
type scheduler struct {
	crawlingSvc     service.CrawlingService
	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

// NewScheduler creates and starts a new job scheduler instance.
func NewScheduler(
	crawlingSvc service.CrawlingService,
	newsSvc service.NewsService,
	systemConfigSvc service.SystemConfigService,
) error {
	svc := &scheduler{
		crawlingSvc:     crawlingSvc,
		newsSvc:         newsSvc,
		systemConfigSvc: systemConfigSvc,
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		return errors.WithStack(err)
	}

	job, err := s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(gocron.NewAtTime(12, 0, 0)),
		),
		gocron.NewTask(svc.crawlingNewsJob),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	logx.Info("newCrawlingNewsJob", job.ID())

	s.Start()

	return nil
}
