//go:build web

package task

import (
	"context"

	"github.com/mjiee/world-news/backend/command"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/tracex"
)

// crawlingNewsJob executes the news crawling job by creating and running a crawling command.
func (s *scheduler) crawlingNewsJob() {
	var (
		ctx = tracex.InjectTraceInContext(context.Background())
		cmd = command.NewCrawlingNewsCommand(ctx, "", nil, nil, s.crawlingSvc, s.newsSvc, s.systemConfigSvc)
	)

	if err := cmd.Execute(ctx); err != nil {
		logx.Error("crawlingNewsJob", err)
	}
}
