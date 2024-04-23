package structs_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trevatk/go-pkg/structs"
)

type DHTNodeSuite struct {
	suite.Suite
	n *structs.DHTNode
}

func TestDHTNodeSuite(t *testing.T) {
	suite.Run(t, new(DHTNodeSuite))
}
