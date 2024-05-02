package dht

import (
	"fmt"
	"net"

	"github.com/structx/go-pkg/domain"
	"github.com/structx/go-pkg/util/encode"
)

// Node kademlia distributed hash table node
type Node struct {
	ID           domain.NodeID
	RoutingTable *domain.RoutingTable
}

// interface compliance
var _ domain.DHT = (*Node)(nil)

// NewNode constructor
func NewNode(ip string, port int) *Node {

	addr := fmt.Sprintf("%s:%d", ip, port)
	nodeID := encode.HashKey([]byte(addr))

	var contactID domain.NodeID = nodeID

	c := &domain.Contact{
		IP:   ip,
		Port: port,
		ID:   contactID,
	}

	rt := &domain.RoutingTable{
		Buckets: map[domain.NodeID][]*domain.Bucket{},
	}

	for i := 0; i < domain.Replication; i++ {
		bucketID := encode.GenerateBucketID(nodeID, i)

		rt.Buckets[nodeID] = append(rt.Buckets[nodeID], &domain.Bucket{
			ID:       bucketID,
			Contacts: make([]*domain.Contact, 0),
		})
	}

	// Insert node contact only in the first bucket
	rt.Buckets[nodeID][0].Contacts = append(rt.Buckets[nodeID][0].Contacts, c)

	return &Node{
		ID:           nodeID,
		RoutingTable: rt,
	}
}

// FindKClosestBuckets iterate over buckets and find shortest distance to key
func (n *Node) FindKClosestBuckets(key []byte) []domain.NodeID {

	hashKey := encode.HashKey(key)
	closestBuckets := make([]domain.NodeID, domain.Replication)

	targetBucketID := encode.GenerateBucketID(hashKey, 0)
	for _, levelBuckets := range n.RoutingTable.Buckets {

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

		for i, bucket := range levelBuckets {
			distance := domain.Distance(targetBucketID).XOR(bucket.ID)
			if compareDistances(distance, bestDistance) < 0 {
				closestBuckets[i] = bucket.ID
			}
		}
	}

	return closestBuckets
}

// FindClosestNodes iterate over bucket and find shortest contact to key
func (n *Node) FindClosestNodes(key, bucketID domain.NodeID) []string {

	closestNodes := make([]string, domain.Replication)

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
			distance := domain.Distance(contact.ID).XOR(key)
			if compareDistances(distance, bestDistance) < 0 {
				closestNodes = append(closestNodes, net.JoinHostPort(contact.IP, fmt.Sprintf("%d", contact.Port)))
			}
		}
	}

	return closestNodes
}

func compareDistances(a, b domain.NodeID) int {
	for i := range a {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
