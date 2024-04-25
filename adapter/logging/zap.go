// Package logging zap logger
package logging

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLoggerFromEnv return uber/zap logger from environment
func NewLoggerFromEnv() (*zap.Logger, error) {

	logFile := os.Getenv("LOG_PATH")
	if logFile == "" {
		return nil, errors.New("$LOG_PATH must be set")
	}

	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if level == "" {
		return nil, errors.New("$LOG_LEVEL must be set")
	}

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		filepath.Clean(logFile),
	}

	var l zapcore.Level
	switch level {
	case "debug":
		l = zap.DebugLevel
	case "error":
		l = zap.ErrorLevel
	}

	cfg.Level = zap.NewAtomicLevelAt(l)

	return cfg.Build()
}
