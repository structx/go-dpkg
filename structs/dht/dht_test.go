package dht_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/structx/go-pkg/structs/dht"
)

type DHTSuite struct {
	suite.Suite
	dht *dht.DHT
}

func (suite *DHTSuite) SetupTest() {

	dht := dht.NewDHT("127.0.0.1:1234")

	suite.dht = dht
}

func (suite *DHTSuite) TestPut() {

	assert := assert.New(suite.T())

	key := []byte("my-key")

	n, err := suite.dht.Put(key)
	assert.NoError(err)

	fmt.Printf("found node %s", hex.EncodeToString(n.ID[:]))
}

func (suite *DHTSuite) TestAddNode() {

	assert := assert.New(suite.T())

	n := dht.NewNode("127.0.0.2:1234")

	assert.NoError(suite.dht.AddNode(n))
}

func TestDHTSuite(t *testing.T) {
	suite.Run(t, new(DHTSuite))
}
