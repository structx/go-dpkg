// Package dht distributed hash table implementation
package dht

import (
	"fmt"
	"net"

	"github.com/structx/go-pkg/domain"
	"github.com/structx/go-pkg/util/encode"
)

// Node kademlia distributed hash table node
type Node struct {
	ID           domain.NodeID224
	RoutingTable *domain.RoutingTable

	replicationFactor int
}

// interface compliance
var _ domain.DHT = (*Node)(nil)

// NewNode constructor
func NewNode(ip string, port, replicationFactor int) *Node {

	addr := fmt.Sprintf("%s:%d", ip, port)
	nodeID := encode.HashKey([]byte(addr))

	var contactID domain.NodeID224 = nodeID

	c := &domain.Contact{
		IP:   ip,
		Port: port,
		ID:   contactID,
	}

	rt := &domain.RoutingTable{
		Buckets: map[domain.NodeID224][]*domain.Bucket{},
	}

	for i := 0; i < replicationFactor; i++ {
		bucketID := encode.MaskFromPrefix(nodeID, i)

		rt.Buckets[nodeID] = append(rt.Buckets[nodeID], &domain.Bucket{
			ID:       bucketID,
			Contacts: make([]*domain.Contact, 0),
		})
	}

	// Insert node contact only in the first bucket
	rt.Buckets[nodeID][0].Contacts = append(rt.Buckets[nodeID][0].Contacts, c)

	return &Node{
		ID:                nodeID,
		RoutingTable:      rt,
		replicationFactor: replicationFactor,
	}
}

// NewNodeWithDefault constructor with default values
func NewNodeWithDefault(ip string, port int) *Node {

	addr := fmt.Sprintf("%s:%d", ip, port)
	nodeID := encode.HashKey([]byte(addr))

	var contactID domain.NodeID224 = nodeID

	c := &domain.Contact{
		Port: port,
		ID:   contactID,
	}
	c.SetID()

	rt := &domain.RoutingTable{
		Buckets: map[domain.NodeID224][]*domain.Bucket{},
	}

	for i := 0; i < domain.DefaultReplicationFactor; i++ {
		bucketID := encode.MaskFromPrefix(nodeID, i)

		rt.Buckets[nodeID] = append(rt.Buckets[nodeID], &domain.Bucket{
			ID:       bucketID,
			Contacts: make([]*domain.Contact, 0),
		})
	}

	// Insert node contact only in the first bucket
	rt.Buckets[nodeID][0].Contacts = append(rt.Buckets[nodeID][0].Contacts, c)

	return &Node{
		ID:                nodeID,
		RoutingTable:      rt,
		replicationFactor: domain.DefaultReplicationFactor,
	}
}

// FindKClosestBuckets iterate over buckets and find shortest distance to key
func (n *Node) FindKClosestBuckets(key []byte) []domain.NodeID224 {

	keyHash := encode.HashKey(key)
	closestBuckets := make([]domain.NodeID224, 0, n.replicationFactor)

	// targetBucketID := encode.MaskFromPrefix(hashKey, 0)
	for nodeID, levelBuckets := range n.RoutingTable.Buckets {

		var bestDistance = [28]byte{
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF,
		}

		for _, bucket := range levelBuckets {
			distance := domain.Distance(keyHash).XOR(bucket.ID)
			if compareDistances(distance, bestDistance) < 0 {
				closestBuckets = append(closestBuckets, nodeID)
			}
		}
	}

	return closestBuckets
}

// FindClosestNodes iterate over bucket and find shortest contact to key
func (n *Node) FindClosestNodes(key []byte, bucketID domain.NodeID224) []string {

	keyHash := encode.HashKey(key)
	closestNodes := make([]string, 0, n.replicationFactor)

	for _, bucket := range n.RoutingTable.Buckets[bucketID] {

		var bestDistance = [28]byte{
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF,
			0xFF,
		}

		for _, contact := range bucket.Contacts {
			distance := domain.Distance(keyHash).XOR(contact.ID)
			if compareDistances(distance, bestDistance) < 0 {
				closestNodes = append(closestNodes, net.JoinHostPort(contact.IP, fmt.Sprintf("%d", contact.Port)))
			}
		}
	}

	return closestNodes
}

// AddOrUpdateNode add or overwrite node
func (n *Node) AddOrUpdateNode(c *domain.Contact) {

	for i := 0; i < n.replicationFactor; i++ {
		bucketID := encode.MaskFromPrefix(c.ID, i)

		n.RoutingTable.Buckets[c.ID] = append(n.RoutingTable.Buckets[c.ID], &domain.Bucket{
			ID:       bucketID,
			Contacts: make([]*domain.Contact, 0),
		})
	}

	// insert contact into bucket
	n.RoutingTable.Buckets[c.ID][0].Contacts = append(n.RoutingTable.Buckets[c.ID][0].Contacts, c)
}

func compareDistances(a, b domain.NodeID224) int {
	for i := range a {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
