package adapter

import (
	"context"

	"github.com/mjiee/world-news/backend/adapter/dto"
	"github.com/mjiee/world-news/backend/command"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/collector"
	"github.com/mjiee/world-news/backend/pkg/databasex"
	"github.com/mjiee/world-news/backend/pkg/httpx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/tracex"
	"github.com/mjiee/world-news/backend/repository"
	"github.com/mjiee/world-news/backend/repository/model"
	"github.com/mjiee/world-news/backend/service"
)

const AppName = "world-news"

// App struct
type App struct {
	ctx context.Context

	crawlingSvc     service.CrawlingService
	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// init is called once at application startup.
func (a *App) init() {
	// init trace
	tracex.InitTracer(AppName)

	// init database
	db, err := databasex.NewAppDB(AppName)
	if err != nil {
		logx.Fatal("database connection failed", err)

		return
	}

	repository.SetDefault(db)

	// auto migrate
	if err := model.AutoMigrate(db); err != nil {
		logx.Fatal("auto migrate database", err)

		return
	}

	// init service
	a.crawlingSvc = service.NewCrawlingService(collector.NewCollector())
	a.newsSvc = service.NewNewsService()
	a.systemConfigSvc = service.NewSystemConfigService()

	// init system config
	if err := a.systemConfigSvc.SystemConfigInit(a.ctx); err != nil {
		logx.Fatal("init system config: %+v", err)
	}
}

// startup is called when the app starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.init()
}

// GetSystemConfig handles the request to retrieve system config.
func (a *App) GetSystemConfig(req *dto.GetSystemConfigRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	data, err := a.systemConfigSvc.GetSystemConfig(ctx, req.Key)

	return httpx.AppRespHandle(ctx, "GetSystemConfig", req, dto.NewSystemConfigFromEntity(data), err)
}

// SaveSystemConfig handles the request to save system config.
func (a *App) SaveSystemConfig(req *dto.SystemConfig) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	return httpx.AppRespHandle(ctx, "SaveSystemConfig", req, nil,
		a.systemConfigSvc.SaveSystemConfig(ctx, req.ToEntity()))
}

// CrawlingNews handles the request to crawl news.
func (a *App) CrawlingNews(req *dto.CrawlingNewsRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	cmd := command.NewCrawlingNewsCommand(a.crawlingSvc, a.newsSvc, a.systemConfigSvc)

	return httpx.AppRespHandle(ctx, "CrawlingNews", req, nil, cmd.Execute(ctx))
}

// CrawlingWebsite handles the request to crawl a news website.
func (a *App) CrawlingWebsite() *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	cmd := command.NewCrawlingNewsWebsiteCommand(a.crawlingSvc, a.systemConfigSvc)

	return httpx.AppRespHandle(ctx, "CrawlingWebsite", nil, nil, cmd.Execute(ctx))
}

// QueryCrawlingRecords handles the request to retrieve crawling records.
func (a *App) QueryCrawlingRecords(req *dto.QueryCrawlingRecordsRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	data, total, err := a.crawlingSvc.QueryCrawlingRecords(ctx,
		*valueobject.NewQueryRecordParams(req.RecordType, req.Status, req.Pagination))

	return httpx.AppRespHandle(ctx, "QueryCrawlingRecords", req, dto.NewQueryCrawlingRecordResult(data, total), err)
}

// DeleteCrawlingRecord handles the request to delete a crawling record.
func (a *App) DeleteCrawlingRecord(req *dto.DeleteCrawlingRecordRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	return httpx.AppRespHandle(ctx, "DeleteCrawlingRecord", req, nil, a.crawlingSvc.DeleteCrawlingRecord(ctx, req.Id))
}

// HasCrawlingTasks handles the request to confirm whether there are ongoing crawling tasks.
func (a *App) HasCrawlingTasks() *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	result, err := a.crawlingSvc.HasProcessingTasks(ctx)

	return httpx.AppRespHandle(ctx, "HasCrawlingTasks", nil, result, err)
}

// QueryNews handles the request to retrieve news detail list.
func (a *App) QueryNews(req *dto.QueryNewsRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	data, total, err := a.newsSvc.QueryNews(ctx,
		valueobject.NewQueryNewsParams(req.RecordId, req.Pagination))

	return httpx.AppRespHandle(ctx, "QueryNews", req, dto.NewQueryNewsResult(data, total), err)
}

// GetNewsDetail handles the request to retrieve a news detail.
func (a *App) GetNewsDetail(req *dto.GetNewsDetailRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	news, err := a.newsSvc.GetNewsDetail(ctx, req.Id)

	return httpx.AppRespHandle(ctx, "GetNewsDetail", req, dto.NewNewsDetailFromEntity(news), err)
}

// DeleteNews handles the request to delete a news detail.
func (a *App) DeleteNews(req *dto.DeleteNewsRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	return httpx.AppRespHandle(ctx, "DeleteNews", req, nil, a.newsSvc.DeleteNews(ctx, req.Id))
}
