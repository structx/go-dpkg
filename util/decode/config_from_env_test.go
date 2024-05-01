package decode_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/structx/go-pkg/adapter/setup"
	"github.com/structx/go-pkg/util/decode"
)

func init() {
	_ = os.Setenv("DSERVICE_CONFIG", "./testfiles/decode.test.hcl")
}

func Test_ConfigFromEnv(t *testing.T) {
	t.Run("simple config", func(t *testing.T) {

		assert := assert.New(t)

		cfg := setup.New()
		assert.NoError(decode.ConfigFromEnv(cfg))

		assert.NotNil(cfg.GetChain())
		assert.NotNil(cfg.GetLogger())
		assert.NotNil(cfg.GetMessenger())
		assert.NotNil(cfg.GetRaft())
		assert.NotNil(cfg.GetServer())
	})
}
