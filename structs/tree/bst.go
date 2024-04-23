// Package tree binary search tree
package tree

// Node tree node
type Node struct {
	payload     int
	left, right *Node
}

// BST binary search tree implementation
type BST struct {
	head *Node
}

// NewBST constructor
func NewBST() *BST {
	return &BST{
		head: nil, // intentially set nil
	}
}

// NewNode constructor
func NewNode(payload int) *Node {
	return &Node{
		payload: payload,
		left:    nil,
		right:   nil,
	}
}

// AddNode to tree
func (b *BST) AddNode(n *Node) {
	if b.head == nil {
		b.head = n
		return
	}

	h := b.head

	if n.payload < b.head.payload {
		if h.left == nil {
			h.left = n
			return
		}
		insertNode(h.left, n)
	} else {
		if h.right == nil {
			h.right = n
			return
		}
		insertNode(h.right, n)
	}
}

func insertNode(n *Node, nn *Node) {

	if nn.payload < n.payload {
		if n.left == nil {
			n.left = nn
			return
		}
		insertNode(n.left, nn)
	} else {
		if n.right == nil {
			n.right = nn
			return
		}
		insertNode(n.right, nn)
	}
}
