// Package logging zap logger
package logging

import (
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/structx/go-pkg/domain"
)

// New uber/zap constructor
func New(cfg domain.Config) (*zap.Logger, error) {

	lcfg := cfg.GetLogger()

	level := strings.ToLower(lcfg.Level)

	zcfg := zap.NewProductionConfig()
	zcfg.OutputPaths = []string{
		filepath.Clean(lcfg.Path),
	}

	var l zapcore.Level
	switch level {
	case "debug":
		l = zap.DebugLevel
	case "error":
		l = zap.ErrorLevel
	}

	zcfg.Level = zap.NewAtomicLevelAt(l)

	return zcfg.Build()
}
