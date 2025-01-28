package adapter

import (
	"context"
	"fmt"

	"github.com/mjiee/world-news/backend/adapter/dto"
	"github.com/mjiee/world-news/backend/pkg/collector"
	"github.com/mjiee/world-news/backend/pkg/databasex"
	"github.com/mjiee/world-news/backend/pkg/httpx"
	"github.com/mjiee/world-news/backend/repository/model"
	"github.com/mjiee/world-news/backend/service"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"gorm.io/gorm"
)

// App struct
type App struct {
	ctx  context.Context
	db   *gorm.DB
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

	a.db = db

	// auto migrate
	if err := model.AutoMigrate(db); err != nil {
		a.logx.Fatal(fmt.Sprintf("Failed to auto migrate database: %+v", err))

		return
	}

	a.newsCrawlingSvc = service.NewNewsCrawlingService(collector.NewCollector(), db)
	a.newsDetailSvc = service.NewNewsDetailService(db)
	a.systemSettingsSvc = service.NewSystemSettingsService(db)
}

// startup is called when the app starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.init()
}

// GetNewsWebsites handles the request to retrieve news websites.
func (a *App) GetNewsWebsites(req *dto.GetNewsWebsitesRequest) *dto.GetNewsWebsitesResponse {
	return nil
}

// AddNewsWebsite handles the request to add a new news website.
func (a *App) AddNewsWebsite(req *dto.AddNewsWebsiteRequest) *httpx.Response {
	return nil
}

// DeleteNewsWebsite handles the request to delete a news website.
func (a *App) DeleteNewsWebsite(req *dto.DeleteNewsWebsiteRequest) *httpx.Response {
	return nil
}

// GetNewsKeywords handles the request to retrieve news keywords.
func (a *App) GetNewsKeywords() *dto.GetNewsKeywordsResponse {
	return nil
}

// DeleteNewsKeyword handles the request to delete a news keyword.
func (a *App) DeleteNewsKeyword(req *dto.DeleteNewsKeywordRequest) *httpx.Response {
	return nil
}

// AddNewsKeyword handles the request to add a new news keyword.
func (a *App) AddNewsKeyword(req *dto.AddNewsKeywordRequest) *httpx.Response {
	return nil
}

// CrawlingNews handles the request to crawl news.
func (a *App) CrawlingNews(req *dto.CrawlingNewsRequest) *httpx.Response {
	return nil
}

// GetCrawlingRecords handles the request to retrieve crawling records.
func (a *App) GetCrawlingRecords(req *dto.GetCrawlingRecordsRequest) *dto.GetCrawlingRecordsResponse {
	return nil
}

// DeleteCrawlingRecord handles the request to delete a crawling record.
func (a *App) DeleteCrawlingRecord(req *dto.DeleteCrawlingRecordRequest) *httpx.Response {
	return nil
}

// GetNewsList handles the request to retrieve news detail list.
func (a *App) GetNewsList(req *dto.GetNewsListRequest) *dto.GetNewsListResponse {
	return nil
}

// GetNewsDetail handles the request to retrieve a news detail.
func (a *App) GetNewsDetail(req *dto.GetNewsDetailRequest) *dto.GetNewsDetailResponse {
	return nil
}

// DeleteNewsDetail handles the request to delete a news detail.
func (a *App) DeleteNewsDetail(req *dto.DeleteNewsRequest) *httpx.Response {
	return nil
}
