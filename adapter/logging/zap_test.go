package logging_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/structx/go-dpkg/adapter/logging"
	"github.com/structx/go-dpkg/adapter/setup"
	"github.com/structx/go-dpkg/util/decode"
)

func init() {
	_ = os.Setenv("DSERVICE_CONFIG", "./testfiles/logger.test.hcl")
}

func Test_NewLoggerFromEnv(t *testing.T) {
	t.Run("provider", func(t *testing.T) {

		assert := assert.New(t)

		cfg := setup.New()
		assert.NoError(decode.ConfigFromEnv(cfg))

		_ = os.Mkdir("./testfiles/log", os.ModePerm)

		logger, err := logging.New(cfg)
		if err != nil {
			t.Fatalf("failed to initialize new logger %v", err)
		}

		logger.Info("success")
		_ = logger.Sync()
	})
}
