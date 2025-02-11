package main

import (
	"fmt"

	"github.com/mjiee/world-news/backend/pkg/databasex"
	"github.com/mjiee/world-news/backend/repository/model"

	"gorm.io/gen"
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

	g.ApplyBasic(model.NewsDetail{}, model.SystemConfig{}, model.CrawlingRecord{})

	g.Execute()
}
