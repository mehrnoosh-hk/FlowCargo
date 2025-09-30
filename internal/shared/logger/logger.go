package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
	out *os.File
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
	var outFile *os.File
	if isDevelopment {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = ""
		}
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
			CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
				var filePath string
				if cwd != "" {
					filePath = strings.TrimPrefix(frame.File, cwd+"/")
				} else {
					filePath = frame.File
				}
				return fmt.Sprintf("%s:%d", filePath, frame.Line), ""
			},
		})
		logger.SetReportCaller(true)
		outFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			logger.SetOutput(outFile)
		} else {
			logger.Info("failed to log to file, using default stderr")
		}
	}
	return Logger{logger, outFile}, nil
}

func (l Logger) Close() error {
	if l.out != nil {
		return l.out.Close()
	}
	return nil
}
