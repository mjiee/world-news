//go:build web

package databasex

import (
	"github.com/mjiee/world-news/backend/pkg/logx"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewWebDB creates a new instance of the gorm.DB
func NewWebDB(addr, logFile string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{
		Logger: logx.NewDBLog(logFile),
	})

	return db, errors.WithStack(err)
}
