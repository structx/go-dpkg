// Package dht distributed hash table implementation
package dht

import (
	"context"
	"fmt"
	"net"

	"github.com/structx/go-dpkg/domain"
	"github.com/structx/go-dpkg/structs/tree"
	"github.com/structx/go-dpkg/util/encode"
)

// Node kademlia distributed hash table node
type Node struct {
	ID           domain.NodeID224
	routingTable *tree.BST

	replicationFactor int
}

// interface compliance
var _ domain.DHT = (*Node)(nil)

// NewNode constructor
func NewNode(ctx context.Context, ip string, port, replicationFactor int) *Node {

	addr := fmt.Sprintf("%s:%d", ip, port)
	nodeID := encode.HashKey([]byte(addr))

	var contactID domain.NodeID224 = nodeID

	c := &domain.Contact{
		IP:   ip,
		Port: port,
		ID:   contactID,
	}

	bst := tree.NewBSTWithDefault()
	bst.Run(ctx)

	bucketSlice := make([]*domain.Bucket, 0, replicationFactor)
	for i := 0; i < replicationFactor; i++ {
		bucketID := encode.MaskFromPrefix(nodeID, i)

		bucketSlice = append(bucketSlice, &domain.Bucket{
			ID:       bucketID,
			Contacts: make([]*domain.Contact, 0),
		})
	}

	// Insert node contact only in the first bucket
	bucketSlice[0].Contacts = append(bucketSlice[0].Contacts, c)

	bst.Insert(context.TODO(), nodeID, bucketSlice)

	return &Node{
		ID:                nodeID,
		routingTable:      bst,
		replicationFactor: replicationFactor,
	}
}

// NewNodeWithDefault constructor with default values
func NewNodeWithDefault(ctx context.Context, ip string, port int) *Node {

	addr := fmt.Sprintf("%s:%d", ip, port)
	nodeID := encode.HashKey([]byte(addr))

	var contactID domain.NodeID224 = nodeID

	c := &domain.Contact{
		Port: port,
		ID:   contactID,
	}
	c.SetID()

	bst := tree.NewBSTWithDefault()
	bst.Run(ctx)

	bucketSlice := make([]*domain.Bucket, 0, domain.DefaultReplicationFactor)
	for i := 0; i < domain.DefaultReplicationFactor; i++ {
		bucketID := encode.MaskFromPrefix(nodeID, i)

		bucketSlice = append(bucketSlice, &domain.Bucket{
			ID:       bucketID,
			Contacts: make([]*domain.Contact, 0),
		})
	}

	// Insert node contact only in the first bucket
	bucketSlice[0].Contacts = append(bucketSlice[0].Contacts, c)

	bst.Insert(context.TODO(), nodeID, bucketSlice)

	return &Node{
		ID:                nodeID,
		routingTable:      tree.NewBSTWithDefault(),
		replicationFactor: domain.DefaultReplicationFactor,
	}
}

// FindKClosestBuckets iterate over buckets and find shortest distance to key
func (n *Node) FindKClosestBuckets(ctx context.Context, key []byte) []domain.NodeID224 {

	keyHash := encode.HashKey(key)

	closestBuckets := make([]domain.NodeID224, 0, n.replicationFactor)

	targetBucketID := encode.MaskFromPrefix(keyHash, 0)
	result := n.routingTable.Search(ctx, targetBucketID)

	if result == nil {
		// no result found
		// include self by default
		return []domain.NodeID224{n.ID}
	}

	bucketSlice, ok := result.GetValue().([]*domain.Bucket)
	if !ok {
		// invalid data structure found
		// include self by default
		return []domain.NodeID224{n.ID}
	}

	// targetBucketID := encode.MaskFromPrefix(hashKey, 0)
	for _, levelBucket := range bucketSlice {

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

		distance := domain.Distance224(keyHash).XOR(levelBucket.ID)
		if compareDistances(distance, bestDistance) < 0 {
			closestBuckets = append(closestBuckets, levelBucket.ID)
		}

	}

	return closestBuckets
}

// FindClosestNodes iterate over bucket and find shortest contact to key
func (n *Node) FindClosestNodes(ctx context.Context, key []byte, bucketID domain.NodeID224) []string {

	keyHash := encode.HashKey(key)
	closestNodes := make([]string, 0, n.replicationFactor)

	result := n.routingTable.Search(ctx, bucketID)
	if result == nil {
		// no bucket was found
		return []string{}
	}

	contactSlice, ok := result.GetValue().([]*domain.Contact)
	if !ok {
		// invalid data structure found
		return []string{}
	}

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

	for _, contact := range contactSlice {
		distance := domain.Distance224(keyHash).XOR(contact.ID)
		if compareDistances(distance, bestDistance) < 0 {
			closestNodes = append(closestNodes, net.JoinHostPort(contact.IP, fmt.Sprintf("%d", contact.Port)))
		}
	}

	return closestNodes
}

// AddOrUpdateNode add or overwrite node
func (n *Node) AddOrUpdateNode(ctx context.Context, c *domain.Contact) {

	for i := 0; i < n.replicationFactor; i++ {
		bucketID := encode.MaskFromPrefix(c.ID, i)

		if i == 0 {
			n.routingTable.Insert(ctx, bucketID, &domain.Bucket{
				ID:       bucketID,
				Contacts: []*domain.Contact{c},
			})
			continue
		}

		n.routingTable.Insert(ctx, bucketID, &domain.Bucket{
			ID:       bucketID,
			Contacts: make([]*domain.Contact, 0),
		})
	}
}

// Get value
func (n *Node) Get(ctx context.Context, key []byte) *domain.Bucket {

	keyHash := encode.HashKey(key)

	result := n.routingTable.Search(ctx, keyHash)
	if result == nil {
		return nil
	}

	bucketSlice, ok := result.GetValue().([]*domain.Bucket)
	if !ok {
		return nil
	}

	for _, levelBucket := range bucketSlice {

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

		for _, contact := range levelBucket.Contacts {
			distance := domain.Distance224(levelBucket.ID).XOR(contact.ID)
			if compareDistances(distance, bestDistance) < 0 {
				return levelBucket
			}
		}
	}

	return nil
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
