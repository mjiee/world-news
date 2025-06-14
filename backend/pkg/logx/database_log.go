package logx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mjiee/world-news/backend/pkg/tracex"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const databaseMsg string = "database"
const slowThreshold = 200

// dbLog is a struct that implements the gorm/logger.Interface interface
type dbLog struct{}

// NewDBLog creates a new instance of the dbLog struct
func NewDBLog(logfile string) *dbLog {
	SetDefaultLogger(logfile)

	return &dbLog{}
}

func (l *dbLog) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *dbLog) Info(ctx context.Context, msg string, data ...interface{}) {
	if defaultLog.Level() > zap.InfoLevel {
		return
	}

	defaultLog.Info(databaseMsg, tracex.LogField(ctx), zap.String(dataField, fmt.Sprintf(msg, data...)))
}

func (l *dbLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	if defaultLog.Level() > zap.WarnLevel {
		return
	}

	defaultLog.Warn(databaseMsg, tracex.LogField(ctx), zap.String(dataField, fmt.Sprintf(msg, data...)))
}

func (l *dbLog) Error(ctx context.Context, msg string, data ...interface{}) {
	if defaultLog.Level() > zap.ErrorLevel {
		return
	}

	defaultLog.Error(databaseMsg, tracex.LogField(ctx), zap.String(dataField, fmt.Sprintf(msg, data...)))
}

func (l *dbLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	var (
		elapsed   = time.Since(begin).Milliseconds()
		sql, rows = fc()
		data      = &LogData{
			Duration: elapsed,
			Rows:     rows,
			SQL:      sql,
		}
	)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		defaultLog.Error(databaseMsg, tracex.LogField(ctx), zap.Any(dataField, data))

		return
	}

	if elapsed > slowThreshold {
		defaultLog.Warn(databaseMsg, tracex.LogField(ctx), zap.Any(dataField, data))

		return
	}

	if defaultLog.Level() > zap.DebugLevel {
		return
	}

	defaultLog.Info(databaseMsg, tracex.LogField(ctx), zap.Any(dataField, data))
}
