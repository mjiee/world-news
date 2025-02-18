package main

import (
	"embed"

	"github.com/mjiee/world-news/cmd"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// run app
	cmd.Run(assets)
}
