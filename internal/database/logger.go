package database

import (
	"context"
	"fmt"
	"time"

	loggerx "github.com/seth16888/wxcommon/logger"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// GormLogger 实现 gorm.Logger 接口
type GormLogger struct {
	LogLevel logger.LogLevel
}

// NewGormLogger 创建新的 GORM 日志适配器
func NewGormLogger(level string) *GormLogger {
	var logLevel logger.LogLevel
	switch level {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	return &GormLogger{
		LogLevel: logLevel,
	}
}

// LogMode 实现 gorm.Logger 接口
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 实现 gorm.Logger 接口
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		loggerx.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn 实现 gorm.Logger 接口
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		loggerx.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error 实现 gorm.Logger 接口
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		loggerx.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace 实现 gorm.Logger 接口
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 根据不同情况记录日志
	switch {
	case err != nil && l.LogLevel >= logger.Error:
		loggerx.Error("GORM SQL Error",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed))
	case elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn:
		loggerx.Warn("GORM Slow SQL",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed))
	case l.LogLevel >= logger.Info:
		loggerx.Info("GORM SQL",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed))
	}
}
