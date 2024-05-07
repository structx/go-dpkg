package setup_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/structx/go-dpkg/adapter/setup"
	"github.com/structx/go-dpkg/util/decode"
)

func init() {
	_ = os.Setenv("DSERVICE_CONFIG", "./testfiles/config.test.hcl")
}

func Test_Provider(t *testing.T) {
	t.Run("provider", func(t *testing.T) {
		assert := assert.New(t)
		cfg := setup.New()
		assert.NoError(decode.ConfigFromEnv(cfg))
	})
}
