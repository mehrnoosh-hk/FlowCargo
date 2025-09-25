package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

func NewLogger(isDevelopment bool, level LogLevel) (Logger, error) {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(string(level))
	if err != nil {
		return Logger{}, err
	}
	logger.SetLevel(lvl)
	if isDevelopment {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: false,
			ForceColors: true,
		})
		logger.SetReportCaller(true)
	}
	return Logger{logger}, nil
}
