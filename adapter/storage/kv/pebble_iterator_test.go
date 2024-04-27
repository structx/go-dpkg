package kv_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trevatk/go-pkg/adapter/logging"
	"github.com/trevatk/go-pkg/adapter/setup"
	"github.com/trevatk/go-pkg/adapter/storage/kv"
	"github.com/trevatk/go-pkg/domain"
	"github.com/trevatk/go-pkg/util/decode"
)

func init() {
	_ = os.Setenv("DSERVICE_CONFIG", "./testfiles/pebble.test.hcl")
}

type PebbleIteratorSuite struct {
	suite.Suite
	it domain.KvIterator
}

func (suite *PebbleIteratorSuite) SetupTest() {
	assert := suite.Assert()

	_ = os.Mkdir("./testfiles/log", os.ModePerm)

	cfg := setup.New()
	assert.NoError(decode.ConfigFromEnv(cfg))

	logger, err := logging.New(cfg)
	assert.NoError(err)

	db, err := kv.NewPebble(logger, cfg)
	assert.NoError(err)

	_ = db.Put([]byte("1"), []byte("1"))
	_ = db.Put([]byte("2"), []byte("2"))
	_ = db.Put([]byte("3"), []byte("3"))

	suite.it, err = db.Iterator(context.TODO())
	assert.NoError(err)
}

func (suite *PebbleIteratorSuite) TestNext() {

	assert := suite.Assert()
	defer func() { assert.NoError(suite.it.Close()) }()

	for suite.it.Next() {
		assert.NotEmpty(suite.it.Key())
	}
}

func TestPebbleIteratorSuite(t *testing.T) {
	suite.Run(t, new(PebbleIteratorSuite))
}
