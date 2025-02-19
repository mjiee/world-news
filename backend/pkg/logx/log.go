package logx

import (
	"context"
	"os"
	"sync"

	"github.com/mjiee/world-news/backend/pkg/tracex"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	dataField = "data"
)

var (
	once       sync.Once
	defaultLog *zap.Logger
)

// LogData is a struct that holds the log data
type LogData struct {
	// common fields
	Duration int64 `json:"duration,omitempty"`

	// http fields
	Method   string `json:"method,omitempty"`
	Request  any    `json:"request,omitempty"`
	Response any    `json:"response,omitempty"`

	// database fields
	SQL  string `json:"sql,omitempty"`
	Rows int64  `json:"rows,omitempty"`
}

// SetDefaultLogger sets the default logger for the application
func SetDefaultLogger(logfile string) {
	once.Do(func() {
		// writer
		var writer zapcore.WriteSyncer

		if logfile != "" {
			logFile := &lumberjack.Logger{
				Filename:   logfile,
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     5,
			}

			writer = zapcore.AddSync(logFile)
		} else {
			writer = zapcore.AddSync(os.Stdout)
		}

		// encoder
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

		// level
		level := zap.NewAtomicLevel()
		level.SetLevel(zapcore.InfoLevel)

		// create logger
		defaultLog = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.NewMultiWriteSyncer(writer),
			level,
		))
	})
}

// checkDefaultLogger checks if the default logger is set
func checkDefaultLogger() {
	if defaultLog == nil {
		SetDefaultLogger("")
	}
}

// info logs an info message
func Info(msg string, data any) {
	checkDefaultLogger()

	defaultLog.Info(msg, zap.Any(dataField, data))
}

// Error logs an error
func Error(msg string, err error) {
	checkDefaultLogger()

	defaultLog.Error(msg, zap.Error(err))
}

// Fatal logs a fatal message
func Fatal(msg string, data error) {
	checkDefaultLogger()

	defaultLog.Fatal(msg, zap.Error(data))
}

// contextLogger is a struct that implements the Logger interface
type contextLogger struct {
	ctx context.Context
}

// WithContext returns a new defaultLogger
func WithContext(ctx context.Context) *contextLogger {
	return &contextLogger{
		ctx: ctx,
	}
}

// info logs an info message
func (l *contextLogger) Info(msg string, data any) {
	checkDefaultLogger()

	defaultLog.Info(msg, tracex.LogField(l.ctx), zap.Any(dataField, data))
}

// Error logs an error
func (l *contextLogger) Error(msg string, err error) {
	checkDefaultLogger()

	defaultLog.Error(msg, tracex.LogField(l.ctx), zap.Error(err))
}
