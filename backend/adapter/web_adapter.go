//go:build web

package adapter

import (
	"context"

	"github.com/mjiee/world-news/backend/adapter/dto"
	"github.com/mjiee/world-news/backend/command"
	"github.com/mjiee/world-news/backend/entity/valueobject"
	"github.com/mjiee/world-news/backend/pkg/collector"
	"github.com/mjiee/world-news/backend/pkg/config"
	"github.com/mjiee/world-news/backend/pkg/databasex"
	"github.com/mjiee/world-news/backend/pkg/httpx"
	"github.com/mjiee/world-news/backend/pkg/tracex"
	"github.com/mjiee/world-news/backend/repository"
	"github.com/mjiee/world-news/backend/repository/model"
	"github.com/mjiee/world-news/backend/service"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const AppName = "world-news"

// WebAadapter struct
type WebAadapter struct {
	crawlingSvc     service.CrawlingService
	newsSvc         service.NewsService
	systemConfigSvc service.SystemConfigService
}

// SetWebAdapter create a new WebAadapter
func SetWebAdapter(conf *config.WebConfig) (*WebAadapter, error) {
	web := &WebAadapter{}

	// init database
	db, err := databasex.NewWebDB(conf.DBAddr, conf.LogFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	repository.SetDefault(db)

	// auto migrate
	if err := model.AutoMigrate(db); err != nil {
		return nil, errors.WithStack(err)
	}

	// init service
	c := collector.NewCollector()

	web.crawlingSvc = service.NewCrawlingService(c)
	web.newsSvc = service.NewNewsService(c)
	web.systemConfigSvc = service.NewSystemConfigService()

	// init system config
	if err := web.systemConfigSvc.SystemConfigInit(context.Background()); err != nil {
		return nil, err
	}

	return web, nil
}

// GetSystemConfig handles the request to retrieve system config.
func (a *WebAadapter) GetSystemConfig(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.GetSystemConfigRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	data, err := a.systemConfigSvc.GetSystemConfig(ctx, req.Key)

	httpx.WebResp(c, dto.NewSystemConfigFromEntity(data), err)
}

// SaveSystemConfig handles the request to save system config.
func (a *WebAadapter) SaveSystemConfig(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.SystemConfig](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	config, err := req.ToEntity()
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	httpx.WebResp(c, nil, a.systemConfigSvc.SaveSystemConfig(ctx, config))
}

// CrawlingNews handles the request to crawling news.
func (a *WebAadapter) CrawlingNews(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.CrawlingNewsRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	cmd := command.NewCrawlingNewsCommand(cmdCtx, req.StartTime, req.Sources, req.Topics,
		a.crawlingSvc, a.newsSvc, a.systemConfigSvc)

	httpx.WebResp(c, nil, cmd.Execute(ctx))
}

// CrawlingWebsite handles the request to crawling website.
func (a *WebAadapter) CrawlingWebsite(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		cmdCtx = tracex.CopyTraceContext(ctx, context.Background())
	)

	cmd := command.NewCrawlingNewsWebsiteCommand(cmdCtx, a.crawlingSvc, a.systemConfigSvc)

	httpx.WebResp(c, nil, cmd.Execute(ctx))
}

// QueryCrawlingRecords handles the request to retrieve crawling records.
func (a *WebAadapter) QueryCrawlingRecords(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.QueryCrawlingRecordsRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	data, total, err := a.crawlingSvc.QueryCrawlingRecords(ctx,
		*valueobject.NewQueryRecordParams(req.RecordType, req.Status, req.Pagination))

	httpx.WebResp(c, dto.NewQueryCrawlingRecordResult(data, total), err)
}

// GetCrawlingRecord handles the request to retrieve a crawling record.
func (a *WebAadapter) GetCrawlingRecord(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.GetCrawlingRecordRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	data, err := a.crawlingSvc.GetCrawlingRecord(ctx, req.Id)

	httpx.WebResp(c, dto.NewCrawlingRecordFromEntity(data), err)
}

// DeleteCrawlingRecord handles the request to delete a crawling record.
func (a *WebAadapter) DeleteCrawlingRecord(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.DeleteCrawlingRecordRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	httpx.WebResp(c, nil, a.crawlingSvc.DeleteCrawlingRecord(ctx, req.Id))
}

// UpdateCrawlingRecordStatus handles the request to update a crawling record status.
func (a *WebAadapter) UpdateCrawlingRecordStatus(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.UpdateCrawlingRecordStatusRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	httpx.WebResp(c, nil, a.crawlingSvc.UpdateCrawlingRecordStatus(ctx, req.Id, req.Status))
}

// HasCrawlingTasks handles the request to confirm whether there are ongoing crawling tasks.
func (a *WebAadapter) HasCrawlingTasks(c *gin.Context) {
	ctx := c.Request.Context()

	result, err := a.crawlingSvc.HasProcessingTasks(ctx)

	httpx.WebResp(c, result, err)
}

// QueryNews handles the request to retrieve news detail list.
func (a *WebAadapter) QueryNews(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.QueryNewsRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	data, total, err := a.newsSvc.QueryNews(ctx, req.ToValueobject())

	httpx.WebResp(c, dto.NewQueryNewsResult(data, total), err)
}

// GetNewsDetail handles the request to retrieve a news detail.
func (a *WebAadapter) GetNewsDetail(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.GetNewsDetailRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	news, err := a.newsSvc.GetNewsDetail(ctx, req.Id)

	httpx.WebResp(c, dto.NewNewsDetailFromEntity(news), err)
}

// DeleteNews handles the request to delete a news detail.
func (a *WebAadapter) DeleteNews(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.DeleteNewsRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	httpx.WebResp(c, nil, a.newsSvc.DeleteNews(ctx, req.Id))
}

// CritiqueNews handles the request to critique a news detail.
func (a *WebAadapter) CritiqueNews(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.CritiqueNewsRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	var (
		cmd       = command.NewCritiqueNewsCommand(req.Title, req.Contents, a.newsSvc, a.systemConfigSvc)
		data, err = cmd.Execute(ctx)
	)

	httpx.WebResp(c, data, err)
}

// TranslateNews handles the request to translate a news detail.
func (a *WebAadapter) TranslateNews(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.TranslateNewsRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	var (
		cmd       = command.NewTranslateNewsCommand(req.Contents, req.ToLang, a.newsSvc, a.systemConfigSvc)
		data, err = cmd.Execute(ctx)
	)

	httpx.WebResp(c, data, err)
}

// SaveNewsFavorite handles the request to save a news favorite.
func (a *WebAadapter) SaveNewsFavorite(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.SaveNewsFavoriteRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	httpx.WebResp(c, nil, a.newsSvc.UpdateNewsFavorite(ctx, req.Id, false))
}

// SaveWebsiteWeight handles the request to save a news website weight.
func (a *WebAadapter) SaveWebsiteWeight(c *gin.Context) {
	ctx, req, err := httpx.ParseRequest[dto.SaveWebsiteWeightRequest](c)
	if err != nil {
		httpx.WebResp(c, nil, err)
		return
	}

	httpx.WebResp(c, nil, a.systemConfigSvc.UpdateNewsWebsiteWeight(ctx, req.Website, req.Step))
}
