package encode_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/structx/go-pkg/domain"
	"github.com/structx/go-pkg/util/encode"
)

func Test_GenerateBucketID(t *testing.T) {
	t.Run("level_zero", func(t *testing.T) {
		hashKey := encode.HashKey([]byte("hello world"))
		nodeID := domain.NodeID(hashKey)
		bucketID := encode.GenerateBucketID(nodeID, 0)
		fmt.Println(hex.Dump(bucketID[:]))
	})
	t.Run("level_one", func(t *testing.T) {
		hashKey := encode.HashKey([]byte("hello world"))
		nodeID := domain.NodeID(hashKey)
		bucketID := encode.GenerateBucketID(nodeID, 1)
		fmt.Println(hex.Dump(bucketID[:]))
	})
}
