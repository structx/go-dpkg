package dht_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/structx/go-pkg/consensus/dht"
)

func Test_NewNode(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		n := dht.NewNode("127.0.0.1", 50051)
		for _, x := range n.RoutingTable.Buckets {
			for _, y := range x {
				fmt.Println(y.ID)
			}
		}
	})
}

func Test_FindKClosestBuckets(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		n := dht.NewNode("127.0.0.1", 50051)
		bucketIDSlice := n.FindKClosestBuckets([]byte("127.0.0.77:50051"))
		for _, bucketID := range bucketIDSlice {
			fmt.Println(hex.EncodeToString(bucketID[:]))
		}
	})
}
