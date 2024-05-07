package tree_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/structx/go-dpkg/structs/tree"
	"github.com/structx/go-dpkg/util/encode"
)

func Test_Insert(t *testing.T) {

	ctx := context.TODO()
	assert := assert.New(t)
	bst := tree.NewBSTWithDefault()

	bst.Run(ctx)
	defer func() { assert.NoError(bst.Close()) }()

	bst.Insert(ctx, encode.HashKey([]byte("6")), "6")
	bst.Insert(ctx, encode.HashKey([]byte("5")), "5")
	bst.Insert(ctx, encode.HashKey([]byte("4")), "4")
}

func Test_Search(t *testing.T) {

	assert := assert.New(t)

	ctx := context.TODO()

	k := encode.HashKey([]byte("1"))
	expected := "1"

	bst := tree.NewBSTWithDefault()
	bst.Run(ctx)
	defer func() { assert.NoError(bst.Close()) }()

	bst.Insert(ctx, encode.HashKey([]byte("3")), "3")
	bst.Insert(ctx, encode.HashKey([]byte("2")), "2")
	bst.Insert(ctx, encode.HashKey([]byte("1")), "1")

	result := bst.Search(ctx, k)
	assert.Equal(expected, result.GetValue())
}

func Test_Delete(t *testing.T) {

	ctx := context.TODO()
	assert := assert.New(t)

	k := encode.HashKey([]byte("2"))
	k1 := encode.HashKey([]byte("1"))

	bst := tree.NewBSTWithDefault()
	bst.Run(ctx)
	defer func() { assert.NoError(bst.Close()) }()

	bst.Insert(ctx, encode.HashKey([]byte("3")), "3")
	bst.Insert(ctx, encode.HashKey([]byte("2")), "2")
	bst.Insert(ctx, encode.HashKey([]byte("1")), "1")

	bst.Delete(ctx, k)

	result := bst.Search(ctx, k1)
	assert.Equal("1", result.GetValue())
}
