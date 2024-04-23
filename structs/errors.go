package structs

import "fmt"

// ErrKeyNotFound
type ErrKeyNotFound struct {
	Key string
}

// String
func (notFound *ErrKeyNotFound) String() string {
	return fmt.Sprintf("key %s not found", notFound.Key)
}
