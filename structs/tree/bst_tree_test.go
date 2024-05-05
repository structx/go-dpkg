package tree_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/structx/go-pkg/structs/tree"
	"github.com/structx/go-pkg/util/encode"
)

func Test_Insert(_ *testing.T) {

	ctx := context.TODO()
	// assert := suite.Assert()
	bst := tree.NewBSTWithDefault()

	bst.Run(ctx)
	defer func() { _ = bst.Close() }()

	op1 := tree.NewInsertParams(encode.HashKey([]byte("6")), "6")
	op2 := tree.NewInsertParams(encode.HashKey([]byte("5")), "5")
	op3 := tree.NewInsertParams(encode.HashKey([]byte("4")), "4")

	bst.ExecuteOp(ctx, op1)
	bst.ExecuteOp(ctx, op2)
	bst.ExecuteOp(ctx, op3)
}

func Test_Search(t *testing.T) {

	assert := assert.New(t)

	ctx := context.TODO()

	k := encode.HashKey([]byte("1"))
	expected := "1"

	bst := tree.NewBSTWithDefault()
	bst.Run(ctx)
	defer func() { _ = bst.Close() }()

	op1 := tree.NewInsertParams(encode.HashKey([]byte("3")), "3")
	op2 := tree.NewInsertParams(encode.HashKey([]byte("2")), "2")
	op3 := tree.NewInsertParams(encode.HashKey([]byte("1")), "1")

	bst.ExecuteOp(ctx, op1)
	bst.ExecuteOp(ctx, op2)
	bst.ExecuteOp(ctx, op3)

	sp := tree.NewSearchParams(k)
	result := bst.ExecuteSearch(ctx, sp)

	v1 := result.(tree.SearchResult)
	assert.Equal(expected, v1.GetValue())
}

func Test_Delete(t *testing.T) {

	ctx := context.TODO()
	assert := assert.New(t)

	k := encode.HashKey([]byte("2"))
	k1 := encode.HashKey([]byte("1"))

	bst := tree.NewBSTWithDefault()
	bst.Run(ctx)
	defer func() { _ = bst.Close() }()

	op1 := tree.NewInsertParams(encode.HashKey([]byte("3")), "3")
	op2 := tree.NewInsertParams(encode.HashKey([]byte("2")), "2")
	op3 := tree.NewInsertParams(encode.HashKey([]byte("1")), "1")

	bst.ExecuteOp(ctx, op1)
	bst.ExecuteOp(ctx, op2)
	bst.ExecuteOp(ctx, op3)

	bst.ExecuteOp(ctx, tree.NewDeleteParams(k))

	output := bst.ExecuteSearch(ctx, tree.NewSearchParams(k1))

	result, _ := output.(tree.SearchResult)
	assert.Equal("1", result.GetValue())
}