package dht_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/structx/go-pkg/domain"
	"github.com/structx/go-pkg/structs/dht"
)

func Test_NewNode(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		n := dht.NewNode("127.0.0.1", 50051, domain.DefaultReplicationFactor)
		for _, x := range n.RoutingTable.Buckets {
			for _, y := range x {
				fmt.Println(y.ID)
			}
		}
	})
}

func Test_FindKClosestBuckets(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		n := dht.NewNode("127.0.0.1", 50051, domain.DefaultReplicationFactor)
		bucketIDSlice := n.FindKClosestBuckets([]byte("127.0.0.77:50051"))
		for _, bucketID := range bucketIDSlice {
			fmt.Println(hex.EncodeToString(bucketID[:]))
		}
	})
}

func Test_FindClosestNodes(t *testing.T) {
	t.Run("default", func(t *testing.T) {

		assert := assert.New(t)

		n := dht.NewNode("127.0.0.1", 50051, domain.DefaultReplicationFactor)
		c := &domain.Contact{
			IP:   "10.0.1.77",
			Port: 50051,
		}
		c.SetID()

		n.AddOrUpdateNode(c)

		k := []byte("10.0.2.35")
		bucketIDSlice := n.FindKClosestBuckets(k)
		addrSlice := n.FindClosestNodes(k, bucketIDSlice[0])
		for _, addr := range addrSlice {
			fmt.Println(addr)
			assert.Equal("10.0.1.77:50051", addr)
		}
	})
}
