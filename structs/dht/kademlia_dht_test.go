package dht_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/structx/go-dpkg/domain"
	"github.com/structx/go-dpkg/structs/dht"
)

func Test_FindKClosestBuckets(t *testing.T) {
	t.Run("default", func(t *testing.T) {

		assert := assert.New(t)
		ctx := context.TODO()
		n := dht.NewNode(ctx, "127.0.0.1", 50051, domain.DefaultReplicationFactor)
		bucketIDSlice := n.FindKClosestBuckets(ctx, []byte("127.0.0.1:50051"))
		for _, bucketID := range bucketIDSlice {
			assert.NotNil(bucketID)
			fmt.Println(hex.EncodeToString(bucketID[:]))
		}
	})
}

func Test_FindClosestNodes(t *testing.T) {
	t.Run("default", func(t *testing.T) {

		assert := assert.New(t)
		ctx := context.TODO()
		n := dht.NewNode(ctx, "127.0.0.1", 50051, domain.DefaultReplicationFactor)
		c := &domain.Contact{
			IP:   "10.0.1.77",
			Port: 50051,
		}
		c.SetID()

		n.AddOrUpdateNode(ctx, c)

		k := []byte("10.0.2.35")
		bucketIDSlice := n.FindKClosestBuckets(ctx, k)
		addrSlice := n.FindClosestNodes(ctx, k, bucketIDSlice[0])
		for _, addr := range addrSlice {
			assert.NotEmpty(addr)
		}
	})
}
