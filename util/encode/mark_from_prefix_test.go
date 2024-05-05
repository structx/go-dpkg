package encode_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/structx/go-pkg/domain"
	"github.com/structx/go-pkg/util/encode"
)

func Test_GenerateBucketID(t *testing.T) {
	t.Run("level_zero", func(t *testing.T) {

		assert := assert.New(t)

		hashKey := encode.HashKey([]byte("hello world"))
		nodeID := domain.NodeID224(hashKey)
		bucketID := encode.MaskFromPrefix(nodeID, 0)
		assert.NotNil(bucketID)
		fmt.Println(hex.Dump(bucketID[:]))
	})
	t.Run("level_one", func(t *testing.T) {

		assert := assert.New(t)

		hashKey := encode.HashKey([]byte("hello world"))
		nodeID := domain.NodeID224(hashKey)
		bucketID := encode.MaskFromPrefix(nodeID, 1)
		assert.NotNil(bucketID)
		fmt.Println(hex.Dump(bucketID[:]))
	})
}
