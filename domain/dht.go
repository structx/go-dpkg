package domain

import (
	"fmt"

	"github.com/structx/go-pkg/util/encode"
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
// 224 bits / 8 bits/byte = 28 bytes
type NodeID224 [28]byte

// Distance between two nodes
type Distance NodeID224

// XOR based Distance between two nodes
func (d Distance) XOR(n NodeID224) NodeID224 {
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
	FindKClosestBuckets(key []byte) []NodeID224
	// FindClosestNodes iterate over buckets and find closest contact addresses
	FindClosestNodes(key []byte, nodeID NodeID224) []string
}
