package logging_test

import (
	"os"
	"testing"

	"github.com/trevatk/go-pkg/adapter/logging"
)

func init() {
	_ = os.Setenv("LOG_LEVEL", "DEBUG")
	_ = os.Setenv("LOG_PATH", "pkg.log")
}

func Test_NewLoggerFromEnv(t *testing.T) {
	t.Run("provider", func(t *testing.T) {
		logger, err := logging.NewLoggerFromEnv()
		if err != nil {
			t.Fatalf("failed to initialize new logger %v", err)
		}

		logger.Info("success")
	})
}
