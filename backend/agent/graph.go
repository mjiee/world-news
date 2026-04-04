package agent

import (
	"context"

	"github.com/mjiee/world-news/backend/service"
)

// GraphService is a service for managing graph
type GraphService struct {
	crawlingSvc     service.CrawlingService
	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

// CreateNewsToPodcastGraph creates a graph for news to podcast
func CreateNewsToPodcastGraph(ctx context.Context) {
	// g := compose.NewGraph[*stage.Stage, *stage.Stage]()
}
