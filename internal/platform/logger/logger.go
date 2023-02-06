package logger

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Logger interface {
	Debug(keyvals ...interface{})
	Info(keyvals ...interface{})
	Warn(keyvals ...interface{})
	Error(keyvals ...interface{})
}

// Log wrapper with two Logger methods
type Log struct {
	Logger log.Logger
}

func (l Log) Debug(keyvals ...interface{}) {
	_ = level.Debug(l.Logger).Log(keyvals...)
}

func (l Log) Info(keyvals ...interface{}) {
	_ = level.Info(l.Logger).Log(keyvals...)
}

func (l Log) Warn(keyvals ...interface{}) {
	_ = level.Warn(l.Logger).Log(keyvals...)
}

func (l Log) Error(keyvals ...interface{}) {
	_ = level.Error(l.Logger).Log(keyvals...)
}

func discoverLogLevel(logLevel string) level.Option {
	switch logLevel {
	case "DEBUG":
		return level.AllowDebug()
	case "INFO":
		return level.AllowInfo()
	case "WARNING":
		return level.AllowWarn()
	case "ERROR":
		return level.AllowError()
	}
	return level.AllowInfo()
}

// InitLogger lazily loads a Logger
func InitLogger(appName, logFormat, logLevel string) Logger {
	var logger log.Logger
	switch logFormat {
	case "json", "JSON":
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	default:
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	}
	logger = level.NewFilter(logger, discoverLogLevel(logLevel))
	logger = log.With(logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.Caller(4),
		"service", appName)
	log := Log{
		Logger: logger,
	}
	return log
}
