//go:build !web

package databasex

import (
	"path/filepath"

	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/pathx"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewAppDB creates a new instance of the gorm.DB
func NewAppDB(appName string) (*gorm.DB, error) {
	dbPath, err := getDbPath(appName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logx.NewDBLog(logx.GetAppLogPath(appName)),
	})
}

// getDbPath returns the path to the SQLite database file
func getDbPath(appName string) (string, error) {
	basePath, err := pathx.GetAppBasePath(appName, "database")
	if err != nil {
		return "", err
	}

	return filepath.Join(basePath, "database.db"), nil
}
