package logger

import (
	"go.uber.org/zap"
)

// Logger is a interface for leveled logging
type Logger interface {
	Log(keyvals ...interface{}) error
	Debug(keyvals ...interface{}) error
	Info(keyvals ...interface{}) error
	Warn(keyvals ...interface{}) error
	Error(keyvals ...interface{}) error
}

// LoggingService is a struct implement logger interface
type LoggingService struct {
	logger *zap.SugaredLogger
}

// NewLogger return a logger instance for app
func NewLogger() *LoggingService {
	return &LoggingService{
		logger: newLogger(),
	}
}

func newLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}
