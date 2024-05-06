// Package queue implementation of queues
package queue

import "sync/atomic"

// Item queue element
type Item struct {
	atomicIndex    atomic.Int64
	atomicPriority atomic.Int64
	value          interface{}
	next           *Item
}

// GetIndex getter item index
func (i *Item) GetIndex() int64 {
	v := i.atomicIndex.Load()
	return v
}

// GetPriority getter item priority
func (i *Item) GetPriority() int64 {
	p := i.atomicPriority.Load()
	return p
}

// GetValue getter item value
func (i *Item) GetValue() interface{} {
	return i.value
}

func newItem(size, priority int64, value interface{}) *Item {

	var i Item

	atomicIndex := atomic.Int64{}
	atomicIndex.Store(priority)

	idx := size + 1

	i.atomicIndex = atomic.Int64{}
	i.atomicIndex.Store(idx)

	i.atomicPriority = atomic.Int64{}
	i.atomicPriority.Store(priority)

	i.next = nil
	i.value = value

	return &i
}

// PriorityQueue implementation
type PriorityQueue struct {
	front *Item
	size  atomic.Int64
}

// NewPriorityQueue constructor
func NewPriorityQueue() *PriorityQueue {

	var pq PriorityQueue

	pq.size = atomic.Int64{}
	pq.size.Store(0)

	pq.front = nil

	return &pq
}

// Pop front item
func (pq *PriorityQueue) Pop() *Item {
	i := pq.front
	pq.front = pq.front.next
	return i
}

// Push item into queue
func (pq *PriorityQueue) Push(priority int64, value interface{}) {

	s := pq.size.Load()
	i := newItem(s, priority, value)

	if pq.front == nil {
		pq.front = i
		return
	}

	current := pq.front
	for current != nil {
		if current.next == nil {
			current.next = i
			return
		}
	}
}

// Peek front item without moving to next
func (pq *PriorityQueue) Peek() *Item {
	return pq.front
}
