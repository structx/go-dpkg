package raftfx_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trevatk/go-pkg/adapter/port/raftfx"
	"github.com/trevatk/go-pkg/adapter/setup"
	"github.com/trevatk/go-pkg/util/decode"
)

func init() {
	_ = os.Setenv("DSERVICE_CONFIG", "./testfiles/raft_config.hcl")
}

func Test_Provider(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		assert := assert.New(t)

		cfg := setup.New()
		err := decode.ConfigFromEnv(cfg)
		assert.NoError(err)

		r, tm, err := raftfx.New(cfg, nil)
		assert.NoError(err)

		assert.NotNil(r)
		assert.NotNil(tm)
	})
}
