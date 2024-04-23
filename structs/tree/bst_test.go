package tree_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trevatk/go-pkg/structs/tree"
)

type BSTSuite struct {
	suite.Suite
	bst *tree.BST
}

func (suite *BSTSuite) SetupTest() {
	suite.bst = tree.NewBST()
}

func (suite *BSTSuite) TestAddNode() {

	n1 := tree.NewNode(3)
	n2 := tree.NewNode(2)
	n3 := tree.NewNode(4)

	suite.bst.AddNode(n1)
	suite.bst.AddNode(n2)
	suite.bst.AddNode(n3)
}

func TestBSTSuite(t *testing.T) {
	suite.Run(t, new(BSTSuite))
}
