package domain

const (
	// Replication defined as maximum bucket size and
	// a replication factor
	//
	// for example, the bucket can only grow to k
	// amount of nodes and it would be unlikely
	// for all nodes in an hour timespan to fail
	Replication = 3
	// Concurrent number of threads to use for connecting to
	// other dht nodes
	Concurrent = 3
)

// NodeID 224 bit sha3 hash
// 224 bits / 8 bits/byte = 28 bytes
type NodeID [28]byte

// Distance between two nodes
type Distance NodeID

// XOR based Distance between two nodes
func (d Distance) XOR(n NodeID) NodeID {
	// perform bitwise XOR operation between node IDs
	result := NodeID{}
	for i := range d {
		result[i] = d[i] ^ n[i]
	}
	return result
}

// Contact to dht node
type Contact struct {
	IP   string
	Port int
	ID   NodeID
}

// Bucket in dht node
type Bucket struct {
	ID       NodeID
	Contacts []*Contact
}

// RoutingTable with logarithmic structure
type RoutingTable struct {
	Buckets map[NodeID][]*Bucket // multi-dimensional slice for buckets at different levels
}

// DHT k-buckets distributed hash table
type DHT interface {
	// FindKClosestBuckets iterate over all buckets and compare key to bucket id
	FindKClosestBuckets(key []byte) []NodeID
	// FindClosestNodes iterate over buckets and find closest contact addresses
	FindClosestNodes(key []byte, nodeID NodeID) []string
}
