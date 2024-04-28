package domain

// Permission application permission enum
type Permission string

const (
	// Read read
	Read Permission = "read"
	// Write write
	Write Permission = "write"
	// Delete delete
	Delete Permission = "delete"
	// Execute execute
	Execute Permission = "execute"
)

// String stringify permission
func (p Permission) String() string {
	return string(p)
}
