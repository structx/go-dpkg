package dht

import "fmt"

// ErrKeyNotFound error
type ErrKeyNotFound struct {
	Key []byte
}

// Error stringify error message
func (notFound *ErrKeyNotFound) Error() string {
	return fmt.Sprintf("key %x not found", notFound.Key)
}
