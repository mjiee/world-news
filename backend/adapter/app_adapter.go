package adapter

import (
	"context"
	"fmt"

	"github.com/mjiee/world-news/backend/adapter/dto"
	"github.com/mjiee/world-news/backend/command"
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

	newsCrawlingSvc   service.NewsCrawlingService
	newsDetailSvc     service.NewsDetailService
	systemSettingsSvc service.SystemSettingsService
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

	a.newsCrawlingSvc = service.NewNewsCrawlingService(collector.NewCollector())
	a.newsDetailSvc = service.NewNewsDetailService()
	a.systemSettingsSvc = service.NewSystemSettingsService()

	if err := a.systemSettingsSvc.SystemConfigInit(a.ctx); err != nil {
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
	return httpx.Resp(a.systemSettingsSvc.GetSystemConfig(a.ctx, req.Key))
}

// SaveSystemConfig handles the request to save system config.
func (a *App) SaveSystemConfig(req *dto.SystemConfig) *httpx.Response {
	return httpx.RespE(a.systemSettingsSvc.SaveSystemConfig(a.ctx, req.ToEntity()))
}

// CrawlingNews handles the request to crawl news.
func (a *App) CrawlingNews(req *dto.CrawlingNewsRequest) *httpx.Response {
	cmd := command.NewCrawlingNewsCommand(a.newsCrawlingSvc, a.newsDetailSvc, a.systemSettingsSvc)

	return httpx.RespE(cmd.Execute(a.ctx))
}

// CrawlingWebsite handles the request to crawl a news website.
func (a *App) CrawlingWebsite() *httpx.Response {
	cmd := command.NewCrawlingNewsWebsiteCommand(a.newsCrawlingSvc, a.systemSettingsSvc)

	return httpx.RespE(cmd.Execute(a.ctx))
}

// QueryCrawlingRecords handles the request to retrieve crawling records.
func (a *App) QueryCrawlingRecords(req *dto.QueryCrawlingRecordsRequest) *httpx.Response {
	data, total, err := a.newsCrawlingSvc.QueryCrawlingRecords(a.ctx, &httpx.Pagination{Page: req.Page, Limit: req.Limit})

	return httpx.Resp(dto.NewQueryCrawlingRecordResult(data, total), err)
}

// DeleteCrawlingRecord handles the request to delete a crawling record.
func (a *App) DeleteCrawlingRecord(req *dto.DeleteCrawlingRecordRequest) *httpx.Response {
	return httpx.RespE(a.newsCrawlingSvc.DeleteCrawlingRecord(a.ctx, req.Id))
}

// QueryNews handles the request to retrieve news detail list.
func (a *App) QueryNews(req *dto.QueryNewsRequest) *httpx.Response {
	data, total, err := a.newsDetailSvc.QueryNews(a.ctx, req.RecordId, req.Pagination)

	return httpx.Resp(dto.NewQueryNewsResult(data, total), err)
}

// GetNewsDetail handles the request to retrieve a news detail.
func (a *App) GetNewsDetail(req *dto.GetNewsDetailRequest) *httpx.Response {
	news, err := a.newsDetailSvc.GetNewsDetail(a.ctx, req.Id)

	return httpx.Resp(dto.NewNewsDetailFromEntity(news), err)
}

// DeleteNewsDetail handles the request to delete a news detail.
func (a *App) DeleteNewsDetail(req *dto.DeleteNewsRequest) *httpx.Response {
	return httpx.RespE(a.newsDetailSvc.DeleteNews(a.ctx, req.Id))
}
