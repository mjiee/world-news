package logx

import (
	"context"
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
	Path     string `json:"path,omitempty"`
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

		logFile := &lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     5,
		}

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

		level := zap.NewAtomicLevel()
		level.SetLevel(zapcore.InfoLevel)

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(logFile)),
			level,
		)

		defaultLog = zap.New(core)
	})
}

// info logs an info message
func Info(msg string, data any) {
	defaultLog.Info(msg, zap.Any(dataField, data))
}

// Error logs an error
func Error(msg string, err error) {
	defaultLog.Error(msg, zap.Error(err))
}

// Fatal logs a fatal message
func Fatal(msg string, data error) {
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
	defaultLog.Info(msg, tracex.LogField(l.ctx), zap.Any(dataField, data))
}

// Error logs an error
func (l *contextLogger) Error(msg string, err error) {
	defaultLog.Error(msg, tracex.LogField(l.ctx), zap.Error(err))
}
