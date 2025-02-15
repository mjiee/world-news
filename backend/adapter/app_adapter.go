package adapter

import (
	"context"
	"fmt"

	"github.com/mjiee/world-news/backend/adapter/dto"
	"github.com/mjiee/world-news/backend/command"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/collector"
	"github.com/mjiee/world-news/backend/pkg/databasex"
	"github.com/mjiee/world-news/backend/pkg/httpx"
	"github.com/mjiee/world-news/backend/repository"
	"github.com/mjiee/world-news/backend/repository/model"
	"github.com/mjiee/world-news/backend/service"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

// App struct
type App struct {
	ctx  context.Context
	logx logger.Logger

	crawlingSvc     service.CrawlingService
	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		logx: logger.NewDefaultLogger(),
	}
}

// init is called once at application startup.
func (a *App) init() {
	db, err := databasex.NewAppDB("world-news")
	if err != nil {
		a.logx.Fatal(fmt.Sprintf("Failed to create database connection: %+v", err))

		return
	}

	repository.SetDefault(db)

	// auto migrate
	if err := model.AutoMigrate(db); err != nil {
		a.logx.Fatal(fmt.Sprintf("Failed to auto migrate database: %+v", err))

		return
	}

	a.crawlingSvc = service.NewCrawlingService(collector.NewCollector())
	a.newsSvc = service.NewNewsService()
	a.systemConfigSvc = service.NewSystemConfigService()

	if err := a.systemConfigSvc.SystemConfigInit(a.ctx); err != nil {
		a.logx.Fatal(fmt.Sprintf("Failed to init system config: %+v", err))
	}
}

// startup is called when the app starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.init()
}

// GetSystemConfig handles the request to retrieve system config.
func (a *App) GetSystemConfig(req *dto.GetSystemConfigRequest) *httpx.Response {
	data, err := a.systemConfigSvc.GetSystemConfig(a.ctx, req.Key)

	return httpx.Resp(dto.NewSystemConfigFromEntity(data), err)
}

// SaveSystemConfig handles the request to save system config.
func (a *App) SaveSystemConfig(req *dto.SystemConfig) *httpx.Response {
	return httpx.RespE(a.systemConfigSvc.SaveSystemConfig(a.ctx, req.ToEntity()))
}

// CrawlingNews handles the request to crawl news.
func (a *App) CrawlingNews(req *dto.CrawlingNewsRequest) *httpx.Response {
	cmd := command.NewCrawlingNewsCommand(a.crawlingSvc, a.newsSvc, a.systemConfigSvc)

	return httpx.RespE(cmd.Execute(a.ctx))
}

// CrawlingWebsite handles the request to crawl a news website.
func (a *App) CrawlingWebsite() *httpx.Response {
	cmd := command.NewCrawlingNewsWebsiteCommand(a.crawlingSvc, a.systemConfigSvc)

	return httpx.RespE(cmd.Execute(a.ctx))
}

// QueryCrawlingRecords handles the request to retrieve crawling records.
func (a *App) QueryCrawlingRecords(req *dto.QueryCrawlingRecordsRequest) *httpx.Response {
	data, total, err := a.crawlingSvc.QueryCrawlingRecords(a.ctx,
		*valueobject.NewQueryRecordParams(req.RecordType, req.Status, req.Pagination))

	return httpx.Resp(dto.NewQueryCrawlingRecordResult(data, total), err)
}

// DeleteCrawlingRecord handles the request to delete a crawling record.
func (a *App) DeleteCrawlingRecord(req *dto.DeleteCrawlingRecordRequest) *httpx.Response {
	return httpx.RespE(a.crawlingSvc.DeleteCrawlingRecord(a.ctx, req.Id))
}

// HasCrawlingTasks handles the request to confirm whether there are ongoing crawling tasks.
func (a *App) HasCrawlingTasks() *httpx.Response {
	return httpx.Resp(a.crawlingSvc.HasProcessingTasks(a.ctx))
}

// QueryNews handles the request to retrieve news detail list.
func (a *App) QueryNews(req *dto.QueryNewsRequest) *httpx.Response {
	data, total, err := a.newsSvc.QueryNews(a.ctx, valueobject.NewQueryNewsParams(req.RecordId, req.Pagination))

	return httpx.Resp(dto.NewQueryNewsResult(data, total), err)
}

// GetNewsDetail handles the request to retrieve a news detail.
func (a *App) GetNewsDetail(req *dto.GetNewsDetailRequest) *httpx.Response {
	news, err := a.newsSvc.GetNewsDetail(a.ctx, req.Id)

	return httpx.Resp(dto.NewNewsDetailFromEntity(news), err)
}

// DeleteNews handles the request to delete a news detail.
func (a *App) DeleteNews(req *dto.DeleteNewsRequest) *httpx.Response {
	return httpx.RespE(a.newsSvc.DeleteNews(a.ctx, req.Id))
}
