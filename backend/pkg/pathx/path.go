package pathx

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

// GetAppBasePath returns the base path for the application based on the operating system
func GetAppBasePath(appName string, subdirectories ...string) (string, error) {
	var basePath string

	switch runtime.GOOS {
	case "darwin": // macOS
		basePath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", appName, filepath.Join(subdirectories...))
	case "linux": // Linux
		basePath = filepath.Join(os.Getenv("HOME"), ".local", "share", appName, filepath.Join(subdirectories...))
	case "windows": // Windows
		basePath = filepath.Join(os.Getenv("APPDATA"), appName, filepath.Join(subdirectories...))
	default:
		return "", errors.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Create the base path directory if it does not exist.
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return "", errors.Errorf("failed to create directory: %v", err)
	}

	return basePath, nil
}
