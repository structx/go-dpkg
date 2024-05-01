package kv_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/structx/go-pkg/adapter/logging"
	"github.com/structx/go-pkg/adapter/setup"
	"github.com/structx/go-pkg/adapter/storage/kv"
	"github.com/structx/go-pkg/domain"
	"github.com/structx/go-pkg/util/decode"
)

func init() {
	_ = os.Setenv("DSERVICE_CONFIG", "./testfiles/pebble.test.hcl")
}

type PebbleDBSuite struct {
	suite.Suite
	db domain.KV
}

func (suite *PebbleDBSuite) SetupTest() {

	assert := suite.Assert()

	_ = os.Mkdir("./testfiles/log", os.ModePerm)

	cfg := setup.New()
	assert.NoError(decode.ConfigFromEnv(cfg))

	logger, err := logging.New(cfg)
	assert.NoError(err)

	db, err := kv.NewPebble(logger, cfg)
	assert.NoError(err)

	suite.db = db
}

func (suite *PebbleDBSuite) TestPut() {

	assert := suite.Assert()

	key := []byte("hello")
	value := []byte("world")

	err := suite.db.Put(key, value)
	assert.NoError(err)

	suite.TeardownTest()
}

func (suite *PebbleDBSuite) TestGet() {

	assert := suite.Assert()

	key := []byte("get")
	value := []byte("value")

	err := suite.db.Put(key, value)
	assert.NoError(err)

	v, err := suite.db.Get(key)
	assert.NoError(err)

	assert.Equal(value, v)

	suite.TeardownTest()
}

func (suite *PebbleDBSuite) TestIterator() {

	assert := suite.Assert()

	_ = suite.db.Put([]byte("1"), []byte("1"))
	_ = suite.db.Put([]byte("2"), []byte("2"))
	_ = suite.db.Put([]byte("3"), []byte("3"))

	it, err := suite.db.Iterator(context.TODO())
	assert.NoError(err)

	for it.Next() {
		assert.NotEmpty(it.Key())
	}
	assert.NoError(it.Close())

	suite.TeardownTest()
}

func (suite *PebbleDBSuite) TeardownTest() {
	assert := suite.Assert()
	assert.NoError(suite.db.Close())
}

func TestPebbleDBSuite(t *testing.T) {
	suite.Run(t, new(PebbleDBSuite))
}
