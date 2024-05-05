package tree_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/structx/go-pkg/structs/tree"
	"github.com/structx/go-pkg/util/encode"
)

func Test_Insert(_ *testing.T) {

	ctx := context.TODO()
	// assert := suite.Assert()
	bst := tree.NewBSTWithDefault()

	go bst.Run(ctx)

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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		op1 := tree.NewInsertParams(encode.HashKey([]byte("3")), "3")
		op2 := tree.NewInsertParams(encode.HashKey([]byte("2")), "2")
		op3 := tree.NewInsertParams(encode.HashKey([]byte("1")), "1")

		bst.ExecuteOp(ctx, op1)
		bst.ExecuteOp(ctx, op2)
		bst.ExecuteOp(ctx, op3)
	}()

	wg.Wait()

	sp := tree.NewSearchParams(k)
	result := bst.ExecuteSearch(ctx, sp)

	v1 := result.(tree.SearchResult)
	assert.Equal(expected, v1.GetValue())
}
