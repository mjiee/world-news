//go:build !web

package cmd

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/mjiee/world-news/backend/adapter"
	"github.com/mjiee/world-news/backend/pkg/config"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/tracex"
)

// Run creates an instance of the desktop application structure and runs it.
func Run(assets embed.FS) {
	// init trace
	tracex.InitTracer(config.AppName)

	// Create an instance of the app structure
	app := adapter.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "WorldNews",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
		Logger: logx.NewAppLog(config.AppName),
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
