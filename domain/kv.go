package domain

// KV key value database interface
type KV interface {
	// Get value by key
	Get(key []byte) ([]byte, error)
	// Put set key/value pair
	Put(key, value []byte) error
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
