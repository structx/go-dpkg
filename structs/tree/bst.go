// Package tree binary search tree implementations
package tree

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/structx/go-dpkg/domain"
	"golang.org/x/sync/errgroup"
)

// opEnum operation type
type opEnum int

const (
	// insert
	insertOp opEnum = iota
	// search
	searchOp
	// delete
	deleteOp
)

// bst operation
type op interface {
	// getter operation type
	getType() opEnum
	// getter operation node key
	getKey() [28]byte
	// getter operation node value
	getValue() interface{}
}

type searchParams struct {
	key [28]byte
	op  opEnum
	ch  chan SearchResult
}

// SearchResult operation search result
type SearchResult struct {
	key   [28]byte
	value interface{}
}

// GetKey getter search result key
func (sr *SearchResult) GetKey() [28]byte {
	return sr.key
}

// GetValue getter search result value
func (sr *SearchResult) GetValue() interface{} {
	return sr.value
}

func (sp *searchParams) getType() opEnum {
	return sp.op
}

func (sp *searchParams) getKey() [28]byte {
	return sp.key
}

// not implemented
func (sp *searchParams) getValue() interface{} {
	return nil
}

type insertParams struct {
	op    opEnum
	key   [28]byte
	value interface{}
	ch    chan struct{}
}

func (ip *insertParams) getType() opEnum {
	return ip.op
}

func (ip *insertParams) getKey() [28]byte {
	return ip.key
}

func (ip *insertParams) getValue() interface{} {
	return ip.value
}

type deleteParams struct {
	op  opEnum
	key [28]byte
	ch  chan struct{}
}

func (dp *deleteParams) getType() opEnum {
	return dp.op
}

func (dp *deleteParams) getKey() [28]byte {
	return dp.key
}

// not implemented
func (dp *deleteParams) getValue() interface{} {
	return nil
}

type node struct {
	atomicKey   atomic.Uint64
	value       interface{}
	left, right *node
}

func newNode(key [28]byte, value interface{}) *node {

	var n node

	v := bytesToUint64(key)

	n.atomicKey = atomic.Uint64{}
	n.atomicKey.Store(v)

	n.value = value

	n.left = nil
	n.right = nil

	return &n
}

// GetKey getter node key
func (n *node) GetKey() [28]byte {
	v := n.atomicKey.Load()
	return uint64ToBytes(v)
}

// GetValue getter node value
func (n *node) GetValue() interface{} {
	return n.value
}

func (n *node) inOrder(key uint64) *node {

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
	head *node

	// concurrency
	cc       int
	errGroup *errgroup.Group
	baseCtx  context.Context
	ch       chan op
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

	b.ch = make(chan op)

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

// Insert into tree
func (b *BST) Insert(ctx context.Context, key [28]byte, value interface{}) {

	p := &insertParams{
		op:    insertOp,
		key:   key,
		value: value,
		ch:    make(chan struct{}),
	}
	var wg sync.WaitGroup
	wg.Add(1)

	// wait for acknowledgement
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-p.ch:
				wg.Done()
				return
			}
		}
	}()

	// send request
	go func() {
		b.ch <- p
	}()

	// wait for acknowledge received
	wg.Wait()
}

// Search on tree
func (b *BST) Search(ctx context.Context, key [28]byte) *SearchResult {

	p := &searchParams{
		key: key,
		op:  searchOp,
		ch:  make(chan SearchResult, 1),
	}
	defer close(p.ch)

	var wg sync.WaitGroup
	wg.Add(1)

	var output SearchResult

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case result, ok := <-p.ch:
				if !ok {
					return
				}
				wg.Done()
				output = result
				return
			}
		}
	}()

	// send request
	go func() {
		b.ch <- p
	}()

	// wait for result
	wg.Wait()

	return &output
}

// Delete key in tree
func (b *BST) Delete(ctx context.Context, key [28]byte) {

	p := &deleteParams{
		key: key,
		op:  deleteOp,
		ch:  make(chan struct{}),
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// wait for acknowledgement
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-p.ch:
				wg.Done()
				return
			}
		}
	}()

	// send request
	go func() {
		b.ch <- p
	}()

	// wait for acknowledge received
	wg.Wait()
}

func (b *BST) insert(n *node) {

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
func (b *BST) traverseInOrder(key [28]byte) *node {

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

func (b *BST) worker(ctx context.Context, ch <-chan op) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		case op, ok := <-ch:
			if !ok {
				return nil
			}

			switch op.getType() {
			case insertOp:

				p, ok := op.(*insertParams)
				if !ok {
					return errors.New("unable to cast operation as insert params")
				}

				n := newNode(p.key, p.value)
				b.insert(n)

				// acknowledge operation completion
				p.ch <- struct{}{}

			case searchOp:
				p, ok := op.(*searchParams)
				if !ok {
					return errors.New("unable to cast operation as search params")
				}

				n := b.traverseInOrder(p.key)
				if n == nil {
					p.ch <- SearchResult{}
					continue
				}

				var v interface{}
				if n.value != nil {
					v = n.value
				}

				p.ch <- SearchResult{
					key:   n.GetKey(),
					value: v,
				}

			case deleteOp:
				p, ok := op.(*deleteParams)
				if !ok {
					return errors.New("unable to cast operation as delete params")
				}

				b.delete(p.key)

				// acknowledge operation completion
				p.ch <- struct{}{}

			default:
				// unsuported operation ignore
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
