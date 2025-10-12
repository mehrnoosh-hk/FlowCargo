package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger interface defines the methods for logging at various levels.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Close() error
}

type loggerImplementation struct {
	*logrus.Logger
	out *os.File
}

// LogLevel is an abstraction over implementations (logrus, zap) log levels.
type LogLevel string

// Log levels supported by the logger.
const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

// NewLogger is a wrapper around the logrus.Logger to implement the Logger interface.
func NewLogger(isDevelopment bool, level LogLevel) (Logger, error) {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(string(level))
	if err != nil {
		return nil, err
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
	return &loggerImplementation{logger, outFile}, nil
}

// Close closes the logger's output file if it was set.
func (l *loggerImplementation) Close() error {
	if l.out != nil {
		return l.out.Close()
	}
	return nil
}
