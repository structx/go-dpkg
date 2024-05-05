// Package tree binary search tree implementations
package tree

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/structx/go-pkg/domain"
	"golang.org/x/sync/errgroup"
)

// OpEnum operation type
type OpEnum int

const (
	// Insert insert
	Insert OpEnum = iota
	// Search search
	Search
	// Delete delete
	Delete
)

// Op operation
type Op interface {
	// GetType getter operation type
	GetType() OpEnum
	// GetKey getter operation node key
	GetKey() [28]byte
	// GetValue getter operation node value
	GetValue() interface{}
}

// SearchParams operation parameters
type SearchParams struct {
	key [28]byte
	op  OpEnum
	ch  chan SearchResult
}

// SearchResult operation result
type SearchResult struct {
	key   [28]byte
	value interface{}
}

// GetKey getter operation result node key
func (sr *SearchResult) GetKey() [28]byte {
	return sr.key
}

// GetValue getter operation result node value
func (sr *SearchResult) GetValue() interface{} {
	return sr.value
}

// NewSearchParams constructor
func NewSearchParams(key [28]byte) *SearchParams {
	return &SearchParams{
		key: key,
		op:  Search,
	}
}

// GetType getter operation type
func (sp *SearchParams) GetType() OpEnum {
	return sp.op
}

// GetKey getter operation node key
func (sp *SearchParams) GetKey() [28]byte {
	return sp.key
}

// GetValue not implemented
func (sp *SearchParams) GetValue() interface{} {
	return nil
}

// InsertParams operation
type InsertParams struct {
	op    OpEnum
	key   [28]byte
	value interface{}
}

// NewInsertParams constructor
func NewInsertParams(key [28]byte, value interface{}) *InsertParams {
	return &InsertParams{
		op:    Insert,
		key:   key,
		value: value,
	}
}

// GetType getter operation type
func (ip *InsertParams) GetType() OpEnum {
	return ip.op
}

// GetKey getter operation node key
func (ip *InsertParams) GetKey() [28]byte {
	return ip.key
}

// GetValue getter operation node value
func (ip *InsertParams) GetValue() interface{} {
	return ip.value
}

// DeleteParams operation
type DeleteParams struct {
	op  OpEnum
	key [28]byte
}

// GetType operation type
func (dp *DeleteParams) GetType() OpEnum {
	return dp.op
}

// GetKey operation node key
func (dp *DeleteParams) GetKey() [28]byte {
	return dp.key
}

// GetValue not implemented
func (dp *DeleteParams) GetValue() interface{} {
	return nil
}

// Node in tree
type Node struct {
	atomicKey   atomic.Uint64
	value       interface{}
	left, right *Node
}

// NewNode constructor
func NewNode(key [28]byte, value interface{}) *Node {

	var n Node

	v := bytesToUint64(key)

	n.atomicKey = atomic.Uint64{}
	n.atomicKey.Store(v)

	n.value = value

	n.left = nil
	n.right = nil

	return &n
}

// GetKey getter node key
func (n *Node) GetKey() [28]byte {
	v := n.atomicKey.Load()
	return uint64ToBytes(v)
}

// GetValue getter node value
func (n *Node) GetValue() interface{} {
	return n.value
}

func (n *Node) inOrder(key uint64) *Node {

	if n == nil {
		return nil
	}

	// check left subtree first
	result := n.left.inOrder(key)
	if result != nil {
		return result
	}

	// check current node key
	a1 := n.atomicKey.Load()
	if a1 == key {
		return n
	}

	// search right subtree after checking left subtree
	return n.right.inOrder(key)
}

// BST binary search tree
type BST struct {
	head *Node

	// concurrency
	cc       int
	errGroup *errgroup.Group
	baseCtx  context.Context
	ch       chan Op
}

// NewBSTWithDefault constructor with default values
func NewBSTWithDefault() *BST {
	return &BST{
		head: nil,
		cc:   domain.Concurrent,
		ch:   nil, // intentionally set nil
	}
}

// Run create and start workers
func (b *BST) Run(ctx context.Context) {

	b.ch = make(chan Op)

	b.errGroup, b.baseCtx = errgroup.WithContext(ctx)
	for i := 0; i < b.cc; i++ {
		b.errGroup.Go(func() error {
			err := b.worker(ctx, b.ch)
			if err != nil {
				return fmt.Errorf("failed worker %v", err)
			}
			return nil
		})
	}

}

// Close channels and wait for err group finish
func (b *BST) Close() error {
	b.baseCtx.Done()
	close(b.ch)
	return b.errGroup.Wait()
}

// ExecuteOp on tree
func (b *BST) ExecuteOp(ctx context.Context, op Op) {
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			b.ch <- op
		}
	}()
}

// ExecuteSearch on tree
func (b *BST) ExecuteSearch(ctx context.Context, op Op) interface{} {

	p, ok := op.(*SearchParams)
	if !ok {
		return nil
	}
	p.ch = make(chan SearchResult, 1)
	defer close(p.ch)

	// send request through channel
	go func() {
		b.ch <- p
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case output, ok := <-p.ch:
			if !ok {
				return nil
			}

			return output
		}
	}
}

func (b *BST) insert(n *Node) {

	if b.head == nil {
		b.head = n
		return
	}

	b1 := n.atomicKey.Load()

	current := b.head
	for {
		a1 := current.atomicKey.Load()

		if b1 < a1 {
			if current.left == nil {
				current.left = n
				return
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = n
				return
			}
			current = current.right
		}
	}
}

// visits nodes in order left, root, right
func (b *BST) traverseInOrder(key [28]byte) *Node {

	keyInt := bytesToUint64(key)

	if b.head == nil {
		return nil
	}

	// check left subtree first
	result := b.head.left.inOrder(keyInt)
	if result != nil {
		return result
	}

	// check root node
	if b.head.atomicKey.Load() == keyInt {
		return b.head
	}

	// search right subtree
	return b.head.right.inOrder(keyInt)
}

func (b *BST) delete(key [28]byte) {

	keyUint := bytesToUint64(key)

	if b.head == nil {
		return
	}

	parent := b.head
	current := b.head

	for current != nil {
		a1 := current.atomicKey.Load()
		if a1 == keyUint {
			break
		} else if keyUint < a1 {
			parent = current
			current = current.left
		} else {
			parent = current
			current = current.right
		}
	}

	if current == nil {
		// node not found
		return
	}

	// remove the node
	if current.left == nil && current.right == nil {
		if parent.left == current {
			parent.left = nil
		} else {
			parent.right = nil
		}
		return
	}

	// node with one child
	if current.left == nil {
		if parent.left == current {
			parent.left = current.right
		} else {
			parent.right = current.right
		}
		return
	}
}

func (b *BST) worker(ctx context.Context, ch <-chan Op) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		case op, ok := <-ch:
			if !ok {
				return errors.New("channel closed")
			}

			switch op.GetType() {
			case Insert:

				p, ok := op.(*InsertParams)
				if !ok {
					return errors.New("unable to cast operation as insert params")
				}

				n := NewNode(p.key, p.value)
				b.insert(n)

			case Search:
				p, ok := op.(*SearchParams)
				if !ok {
					return errors.New("unable to cast operation as search params")
				}

				n := b.traverseInOrder(p.key)
				if n == nil || n.value == nil {
					p.ch <- SearchResult{}
					continue
				}
				p.ch <- SearchResult{
					key:   n.GetKey(),
					value: n.value,
				}

			case Delete:
				p, ok := op.(*DeleteParams)
				if !ok {
					return errors.New("unable to cast operation as delete params")
				}

				b.delete(p.key)
			default:
				// unsuported operation
				continue
			}
		}
	}

}

func bytesToUint64(key [28]byte) uint64 {
	var result uint64
	for i, b := range key {
		result ^= uint64(b) << (uint(i) * 8)
	}
	return result
}

func uint64ToBytes(value uint64) [28]byte {
	var result [28]byte
	for i := range result {
		result[i] = byte(value >> (uint(i) * 8))
	}
	return result
}
