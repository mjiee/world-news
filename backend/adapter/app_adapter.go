//go:build !web

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
	ctx    context.Context
	cancel context.CancelFunc

	crawlingSvc     service.CrawlingService
	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{}

	// init database
	db, err := databasex.NewAppDB(AppName)
	if err != nil {
		logx.Fatal("NewAppDB", err)

		return app
	}

	repository.SetDefault(db)

	// auto migrate
	if err := model.AutoMigrate(db); err != nil {
		logx.Fatal("AutoMigrate", err)

		return app
	}

	// init service
	app.crawlingSvc = service.NewCrawlingService(collector.NewCollector())
	app.newsSvc = service.NewNewsService()
	app.systemConfigSvc = service.NewSystemConfigService()

	return app
}

// startup is called when the app starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)

	// init system config
	if err := a.systemConfigSvc.SystemConfigInit(a.ctx); err != nil {
		logx.Fatal("SystemConfigInit", err)
	}
}

// Shutdown is called at application termination.
func (a *App) Shutdown(ctx context.Context) {
	a.cancel()

	if err := a.crawlingSvc.PauseAllTasks(ctx); err != nil {
		logx.Fatal("PauseAllTasks", err)
	}
}

// GetSystemConfig handles the request to retrieve system config.
func (a *App) GetSystemConfig(req *dto.GetSystemConfigRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	data, err := a.systemConfigSvc.GetSystemConfig(ctx, req.Key)

	return httpx.AppResp(ctx, "GetSystemConfig", req, dto.NewSystemConfigFromEntity(data), err)
}

// SaveSystemConfig handles the request to save system config.
func (a *App) SaveSystemConfig(req *dto.SystemConfig) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	return httpx.AppResp(ctx, "SaveSystemConfig", req, nil,
		a.systemConfigSvc.SaveSystemConfig(ctx, req.ToEntity()))
}

// CrawlingNews handles the request to crawl news.
func (a *App) CrawlingNews(req *dto.CrawlingNewsRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	cmd := command.NewCrawlingNewsCommand(a.ctx, req.StartTime, a.crawlingSvc, a.newsSvc, a.systemConfigSvc)

	return httpx.AppResp(ctx, "CrawlingNews", req, nil, cmd.Execute(ctx))
}

// CrawlingWebsite handles the request to crawl a news website.
func (a *App) CrawlingWebsite() *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	cmd := command.NewCrawlingNewsWebsiteCommand(a.ctx, a.crawlingSvc, a.systemConfigSvc)

	return httpx.AppResp(ctx, "CrawlingWebsite", nil, nil, cmd.Execute(ctx))
}

// QueryCrawlingRecords handles the request to retrieve crawling records.
func (a *App) QueryCrawlingRecords(req *dto.QueryCrawlingRecordsRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	data, total, err := a.crawlingSvc.QueryCrawlingRecords(ctx,
		*valueobject.NewQueryRecordParams(req.RecordType, req.Status, req.Pagination))

	return httpx.AppResp(ctx, "QueryCrawlingRecords", req, dto.NewQueryCrawlingRecordResult(data, total), err)
}

// GetCrawlingRecord handles the request to retrieve a crawling record.
func (a *App) GetCrawlingRecord(req *dto.GetCrawlingRecordRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	data, err := a.crawlingSvc.GetCrawlingRecord(ctx, req.Id)

	return httpx.AppResp(ctx, "GetCrawlingRecord", req, dto.NewCrawlingRecordFromEntity(data), err)
}

// DeleteCrawlingRecord handles the request to delete a crawling record.
func (a *App) DeleteCrawlingRecord(req *dto.DeleteCrawlingRecordRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	return httpx.AppResp(ctx, "DeleteCrawlingRecord", req, nil, a.crawlingSvc.DeleteCrawlingRecord(ctx, req.Id))
}

// UpdateCrawlingRecordStatus handles the request to update a crawling record status.
func (a *App) UpdateCrawlingRecordStatus(req *dto.UpdateCrawlingRecordStatusRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	return httpx.AppResp(ctx, "UpdateCrawlingRecordStatus", req, nil,
		a.crawlingSvc.UpdateCrawlingRecordStatus(ctx, req.Id, req.Status))
}

// HasCrawlingTasks handles the request to confirm whether there are ongoing crawling tasks.
func (a *App) HasCrawlingTasks() *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	result, err := a.crawlingSvc.HasProcessingTasks(ctx)

	return httpx.AppResp(ctx, "HasCrawlingTasks", nil, result, err)
}

// QueryNews handles the request to retrieve news detail list.
func (a *App) QueryNews(req *dto.QueryNewsRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	data, total, err := a.newsSvc.QueryNews(ctx, req.ToValueobject())

	return httpx.AppResp(ctx, "QueryNews", req, dto.NewQueryNewsResult(data, total), err)
}

// GetNewsDetail handles the request to retrieve a news detail.
func (a *App) GetNewsDetail(req *dto.GetNewsDetailRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	news, err := a.newsSvc.GetNewsDetail(ctx, req.Id)

	return httpx.AppResp(ctx, "GetNewsDetail", req, dto.NewNewsDetailFromEntity(news), err)
}

// DeleteNews handles the request to delete a news detail.
func (a *App) DeleteNews(req *dto.DeleteNewsRequest) *httpx.Response {
	ctx := tracex.InjectTraceInContext(a.ctx)

	return httpx.AppResp(ctx, "DeleteNews", req, nil, a.newsSvc.DeleteNews(ctx, req.Id))
}

// CritiqueNews handles the request to critique a news detail.
func (a *App) CritiqueNews(req *dto.CritiqueNewsRequest) *httpx.Response {
	var (
		ctx = tracex.InjectTraceInContext(a.ctx)
		cmd = command.NewCritiqueNewsCommand(req.Id, a.newsSvc, a.systemConfigSvc)
	)

	data, err := cmd.Execute(ctx)

	return httpx.AppResp(ctx, "CritiqueNews", req, data, err)
}

// TranslateNews handles the request to translate a news detail.
func (a *App) TranslateNews(req *dto.TranslateNewsRequest) *httpx.Response {
	var (
		ctx = tracex.InjectTraceInContext(a.ctx)
		cmd = command.NewTranslateNewsCommand(req.Id, req.Texts, req.ToLang, a.newsSvc, a.systemConfigSvc)
	)

	data, err := cmd.Execute(ctx)

	return httpx.AppResp(ctx, "TranslateNews", req, data, err)
}
