package domain

import (
	"context"
	"fmt"

	"github.com/structx/go-dpkg/util/encode"
)

const (
	// DefaultReplicationFactor defined as maximum bucket size and
	// a replication factor
	//
	// for example, the bucket can only grow to k
	// amount of nodes and it would be unlikely
	// for all nodes in an hour timespan to fail
	DefaultReplicationFactor = 3
	// Concurrent number of threads to use for connecting to
	// other dht nodes
	Concurrent = 3
)

// NodeID224 224 bit sha3 hash
type NodeID224 [28]byte // 224 bits / 8 bits/byte = 28 bytes
// NodeID256 256 bit sha3 hash
type NodeID256 [32]byte // 256 bits / 8 bits/byte = 32 bytes
// NodeID384 384 bit sha3 hash
type NodeID384 [48]byte // 384 bits / 8 bits/byte = 48 bytes
// NodeID512 512 bit sha3 hash
type NodeID512 [64]byte // 512 bits / 8 bits/byte = 64 bytes

// Distance224 between two nodes
type Distance224 NodeID224

// XOR based Distance between two nodes
func (d Distance224) XOR(n NodeID224) NodeID224 {
	// perform bitwise XOR operation between node IDs
	result := NodeID224{}
	for i := range d {
		result[i] = d[i] ^ n[i]
	}
	return result
}

// Contact to dht node
type Contact struct {
	IP   string
	Port int
	ID   NodeID224
}

// SetID of contact
func (c *Contact) SetID() {
	c.ID = encode.HashKey([]byte(fmt.Sprintf("%s:%d", c.IP, c.Port)))
}

// Bucket in dht node
type Bucket struct {
	ID       NodeID224
	Contacts []*Contact
}

// RoutingTable with logarithmic structure
type RoutingTable struct {
	Buckets map[NodeID224][]*Bucket // multi-dimensional slice for buckets at different levels
}

// DHT k-buckets distributed hash table
type DHT interface {
	// FindKClosestBuckets iterate over all buckets and compare key to bucket id
	FindKClosestBuckets(ctx context.Context, key []byte) []NodeID224
	// FindClosestNodes iterate over buckets and find closest contact addresses
	FindClosestNodes(ctx context.Context, key []byte, nodeID NodeID224) []string
	// AddOrUpdateNode add or override node value
	AddOrUpdateNode(ctx context.Context, c *Contact)
	// Get value
	Get(ctx context.Context, key []byte) *Bucket
}
