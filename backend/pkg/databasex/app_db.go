package databasex

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewAppDB creates a new instance of the AppDB struct
func NewAppDB(appName string) (*gorm.DB, error) {
	dbPath, err := getDbPath(appName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}

// getDbPath returns the file path for the database based on the application name and the operating system
func getDbPath(appName string) (string, error) {
	var basePath string

	switch runtime.GOOS {
	case "darwin": // macOS
		basePath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", appName)
	case "linux": // Linux
		basePath = filepath.Join(os.Getenv("HOME"), ".local", "share", appName)
	case "windows": // Windows
		basePath = filepath.Join(os.Getenv("APPDATA"), appName)
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Create the base path directory if it does not exist.
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	return filepath.Join(basePath, "database.db"), nil
}
