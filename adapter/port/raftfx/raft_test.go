package raftfx_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/structx/go-dpkg/adapter/port/raftfx"
	"github.com/structx/go-dpkg/adapter/setup"
	"github.com/structx/go-dpkg/util/decode"
)

func init() {
	_ = os.Setenv("DSERVICE_CONFIG", "./testfiles/raft.test.hcl")
}

func Test_New(t *testing.T) {
	t.Run("provider", func(t *testing.T) {

		_ = os.Mkdir("./testfiles/log", os.ModePerm)
		_ = os.Mkdir("./testfiles/raft", os.ModePerm)

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
