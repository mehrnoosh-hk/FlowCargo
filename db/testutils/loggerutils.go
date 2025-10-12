package testutils

import (
	"fmt"

	"flowcargo/internal/shared/logger"
)

// TestLogger is a custom logger for testing purposes that captures log messages.
type TestLogger struct {
	DebugMessages []string
	InfoMessages  []string
	WarnMessages  []string
	ErrorMessages []string
	FatalMessages []string
	IsClosed      bool
}

// NewTestLogger creates a new TestLogger instance for capturing log messages during tests.
func NewTestLogger() logger.Logger {
	return &TestLogger{
		DebugMessages: []string{},
		InfoMessages:  []string{},
		WarnMessages:  []string{},
		ErrorMessages: []string{},
		FatalMessages: []string{},
		IsClosed:      false,
	}
}

// Close marks the logger as closed.
func (tl *TestLogger) Close() error {
	tl.IsClosed = true
	return nil
}

// Debug captures debug messages.
func (tl *TestLogger) Debug(args ...interface{}) {
	tl.DebugMessages = append(tl.DebugMessages, fmt.Sprint(args...))
}

// Debugf captures formatted debug messages.
func (tl *TestLogger) Debugf(format string, args ...interface{}) {
	tl.DebugMessages = append(tl.DebugMessages, fmt.Sprintf(format, args...))
}

// Info captures info messages.
func (tl *TestLogger) Info(args ...interface{}) {
	tl.InfoMessages = append(tl.InfoMessages, fmt.Sprint(args...))
}

// Infof captures formatted info messages.
func (tl *TestLogger) Infof(format string, args ...interface{}) {
	tl.InfoMessages = append(tl.InfoMessages, fmt.Sprintf(format, args...))
}

// Warn captures warning messages.
func (tl *TestLogger) Warn(args ...interface{}) {
	tl.WarnMessages = append(tl.WarnMessages, fmt.Sprint(args...))
}

// Warnf captures formatted warning messages.
func (tl *TestLogger) Warnf(format string, args ...interface{}) {
	tl.WarnMessages = append(tl.WarnMessages, fmt.Sprintf(format, args...))
}

// Error captures error messages.
func (tl *TestLogger) Error(args ...interface{}) {
	tl.ErrorMessages = append(tl.ErrorMessages, fmt.Sprint(args...))
}

// Errorf captures formatted error messages.
func (tl *TestLogger) Errorf(format string, args ...interface{}) {
	tl.ErrorMessages = append(tl.ErrorMessages, fmt.Sprintf(format, args...))
}

// Fatal captures fatal messages.
func (tl *TestLogger) Fatal(args ...interface{}) {
	tl.FatalMessages = append(tl.FatalMessages, fmt.Sprint(args...))
}

// Fatalf captures formatted fatal messages.
func (tl *TestLogger) Fatalf(format string, args ...interface{}) {
	tl.FatalMessages = append(tl.FatalMessages, fmt.Sprintf(format, args...))
}
