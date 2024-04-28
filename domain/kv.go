package domain

import "context"

// KV key value database interface
//
//go:generate mockery --name KV
type KV interface {
	// Get value by key
	Get(key []byte) ([]byte, error)
	// Put set key/value pair
	Put(key, value []byte) error
	// Iterator key/value iterator
	Iterator(ctx context.Context) (KvIterator, error)
	// Close database connection
	Close() error
}

// KvIterator key value database interface
type KvIterator interface {
	// Next if next keyvalue pair is not null
	Next() bool
	// Key getter key from current index
	Key() []byte
	// Close iterator
	Close() error
}
