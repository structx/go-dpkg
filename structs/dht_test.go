package structs_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trevatk/go-pkg/structs"
)

type DHTSuite struct {
	suite.Suite
	dht *structs.DHT
}

func (suite *DHTSuite) SetupTest() {

	assert := assert.New(suite.T())

	dht, err := structs.NewDHT("127.0.0.1:1234")
	assert.NoError(err)

	suite.dht = dht
}

func (suite *DHTSuite) TestPut() {

	assert := assert.New(suite.T())

	data := "some data"
	key := []byte("my-key")

	n, err := suite.dht.Put(key, []byte(data))
	assert.NoError(err)

	fmt.Printf("found node %s", hex.EncodeToString(n.ID[:]))
}

func TestDHTSuite(t *testing.T) {
	suite.Run(t, new(DHTSuite))
}
