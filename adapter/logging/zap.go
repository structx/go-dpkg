// Package logging zap logger
package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/structx/go-dpkg/domain"
)

var (
	c zapcore.Core
)

// New uber/zap constructor
func New(cfg domain.Config) (*zap.Logger, error) {

	lcfg := cfg.GetLogger()
	level := strings.ToLower(lcfg.Level)

	f, err := os.OpenFile(filepath.Clean(lcfg.Path), os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s %v", lcfg.Path, err)
	}

	zcfg := zap.NewProductionConfig()

	var l zapcore.Level
	switch level {
	case "debug":
		l = zap.DebugLevel
	case "info":
		l = zap.InfoLevel
	case "warn":
		l = zap.WarnLevel
	case "error":
		// used for production
		l = zap.ErrorLevel
		zcfg.Development = false
	}
	zcfg.Level = zap.NewAtomicLevelAt(l)

	// write syncers
	stdoutSyncer := zapcore.Lock(os.Stdout)
	stderrSyncer := zapcore.Lock(os.Stderr)

	c = zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zcfg.EncoderConfig),
			stdoutSyncer,
			l,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zcfg.EncoderConfig),
			stderrSyncer,
			zap.ErrorLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(f),
			l,
		),
	)

	return zcfg.Build(zap.WrapCore(zapCore))
}

func zapCore(_ zapcore.Core) zapcore.Core {
	return c
}
