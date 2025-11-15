//go:build !web

package logx

import (
	"fmt"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/mjiee/world-news/backend/pkg/pathx"
)

const appMsg string = "app"

// appLog is a struct that implements the wails/logger.Logger interface
type appLog struct{}

// NewAppLog creates a new instance of the appLog struct
func NewAppLog(appName string) *appLog {
	SetDefaultLogger(GetAppLogPath(appName))

	return &appLog{}
}

// GetAppLogPath returns the path to the log file for the application
func GetAppLogPath(appName string) string {
	basePath, err := pathx.GetAppBasePath(appName, "logs")
	if err != nil {
		panic(err)
	}

	return filepath.Join(basePath, fmt.Sprintf("%s.log", appName))
}

func (a *appLog) Print(message string) { defaultLog.Info(appMsg, zap.String(dataField, message)) }

func (a *appLog) Trace(message string) { defaultLog.Info(appMsg, zap.String(dataField, message)) }

func (a *appLog) Debug(message string) { defaultLog.Debug(appMsg, zap.String(dataField, message)) }

func (a *appLog) Info(message string) { defaultLog.Info(appMsg, zap.String(dataField, message)) }

func (a *appLog) Warning(message string) { defaultLog.Warn(appMsg, zap.String(dataField, message)) }

func (a *appLog) Error(message string) { defaultLog.Error(appMsg, zap.String(dataField, message)) }

func (a *appLog) Fatal(message string) { defaultLog.Fatal(appMsg, zap.String(dataField, message)) }
