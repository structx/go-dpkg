package kv

import (
	"context"
	"fmt"

	"github.com/cockroachdb/pebble"
	"github.com/trevatk/go-pkg/domain"
)

// PebbleIterator kv iterator implementation
type PebbleIterator struct {
	it *pebble.Iterator
}

// interface compliance
var _ domain.KvIterator = (*PebbleIterator)(nil)

// Iterator iterator constructor
func (p *PebbleDB) Iterator(ctx context.Context) (domain.KvIterator, error) {

	opts := &pebble.IterOptions{}
	it, err := p.db.NewIterWithContext(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize iterator %v", err)
	}

	return &PebbleIterator{
		it: it,
	}, nil
}

// Next if next keyvalue pair is not null
func (pi *PebbleIterator) Next() bool {
	return pi.it.Next()
}

// Key getter key from current index
func (pi *PebbleIterator) Key() []byte {
	return pi.it.Key()
}

// Close iterator
func (pi *PebbleIterator) Close() error {
	return pi.it.Close()
}
