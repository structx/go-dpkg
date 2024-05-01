package serverfx_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/structx/go-pkg/adapter/port/http/serverfx"
	"github.com/structx/go-pkg/adapter/setup"
	"github.com/structx/go-pkg/util/decode"
)

func init() {
	_ = os.Setenv("DSERVICE_CONFIG", "./testfiles/server.test.hcl")
}

func Test_New(t *testing.T) {
	t.Run("provider", func(t *testing.T) {
		assert := assert.New(t)

		cfg := setup.New()
		assert.NoError(decode.ConfigFromEnv(cfg))

		s := serverfx.New(cfg, nil)

		assert.NotNil(s)
	})
}
