// Package kv database implementation
package kv

import (
	"fmt"
	"path/filepath"

	"github.com/cockroachdb/pebble"

	"go.uber.org/zap"

	"github.com/trevatk/go-pkg/domain"
)

// PebbleDB kv database implementation
// serves as wrapper around cockroad pebble db
type PebbleDB struct {
	db *pebble.DB
}

// interface compliance
var _ domain.KV = (*PebbleDB)(nil)

// NewPebble return new pebble db wrapper class
func NewPebble(logger *zap.Logger, cfg domain.Config) (*PebbleDB, error) {

	ccfg := cfg.GetChain()
	suggaredLogger := logger.Named("PebbleRepository").Sugar()

	opts := &pebble.Options{Logger: suggaredLogger}
	db, err := pebble.Open(filepath.Clean(ccfg.BaseDir), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open pebble db: %v", err)
	}

	return &PebbleDB{
		db: db,
	}, nil
}

// Put set key/value pair
func (p *PebbleDB) Put(key, value []byte) error {
	return p.db.Set(key, value, pebble.Sync)
}

// Get value by key
func (p *PebbleDB) Get(key []byte) ([]byte, error) {

	v, closer, err := p.db.Get(key)
	if err != nil && err == pebble.ErrNotFound {
		return []byte{}, &ErrNotFound{Key: key}
	} else if err != nil {
		return []byte{}, fmt.Errorf("failed to get key value %v", err)
	}

	err = closer.Close()
	if err != nil {
		return []byte{}, fmt.Errorf("failed to close closer %v", err)
	}

	return v, nil
}

// Close database connection
func (p *PebbleDB) Close() error {
	return p.db.Close()
}
