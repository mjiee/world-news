//go:build web

package cmd

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/mjiee/world-news/backend/adapter"
	"github.com/mjiee/world-news/backend/pkg/config"
	"github.com/mjiee/world-news/backend/pkg/locale"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/tracex"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Run creates an instance of the web application structure and runs it.
func Run(assets embed.FS) {
	// init trace
	tracex.InitTracer(adapter.AppName)

	// init config
	config := config.NewWebConfig()

	// set default loger
	logx.SetDefaultLogger(config.LogFile)

	// web adapter
	webAdapter, err := adapter.SetWebAdapter(config)
	if err != nil {
		logx.Fatal("init web adapter", err)
	}

	// service
	r := gin.New()

	// middleware
	r.Use(tracex.WebTracer()).
		Use(logx.WebLogger()).
		Use(locale.WebLocale()).
		Use(cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"POST", "GET", "OPTIONS"},
			AllowHeaders:    []string{"Origin", "Content-Type", "Content-Length", "Content-Language", "Authorization"},
		}))

	// register api router
	ApiRouter(r.Group("/api"), webAdapter)

	// serve static files
	staticFp, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		logx.Fatal("init static file", err)
	}

	r.StaticFS("/", http.FS(staticFp))

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.AbortWithStatus(http.StatusNotFound)

			return
		}

		c.Redirect(http.StatusMovedPermanently, "/")
	})

	// run app
	if err = r.Run(config.Host); err != nil {
		logx.Fatal("run web app", err)
	}
}

// ApiRouter registers the API routes.
func ApiRouter(r *gin.RouterGroup, webAdapter *adapter.WebAadapter) {
	r.POST("/system/config", webAdapter.GetSystemConfig)
	r.POST("/system/config/save", webAdapter.SaveSystemConfig)
	r.POST("/crawling/website", webAdapter.CrawlingWebsite)
	r.POST("/crawling/news", webAdapter.CrawlingNews)
	r.POST("/crawling/processing/task", webAdapter.HasCrawlingTasks)
	r.POST("/crawling/record/detail", webAdapter.GetCrawlingRecord)
	r.POST("/crawling/record/query", webAdapter.QueryCrawlingRecords)
	r.POST("/crawling/record/delete", webAdapter.DeleteCrawlingRecord)
	r.POST("/crawling/record/status", webAdapter.UpdateCrawlingRecordStatus)
	r.POST("/news/query", webAdapter.QueryNews)
	r.POST("/news/detail", webAdapter.GetNewsDetail)
	r.POST("/news/delete", webAdapter.DeleteNews)
}
