package main

import (
	"fmt"

	"gorm.io/gen"

	"github.com/mjiee/world-news/backend/pkg/databasex"
	"github.com/mjiee/world-news/backend/repository/model"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "backend/repository",
		Mode:    gen.WithDefaultQuery,
	})

	db, err := databasex.NewAppDB("world-news")
	if err != nil {
		fmt.Printf("%+v", err)

		return
	}

	g.UseDB(db)

	g.ApplyBasic(model.NewsDetail{}, model.SystemConfig{}, model.CrawlingRecord{}, model.PodcastTask{}, model.Podcast{})

	g.Execute()
}
