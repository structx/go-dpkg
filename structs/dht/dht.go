// Package dht distributed hash table
package dht

import (
	"fmt"

	"golang.org/x/crypto/sha3"
)

var (
	k = 3
)

// Node distributed hash table node
type Node struct {
	ID           [28]byte
	Addr         string
	routingTable map[[28]byte]*Node
}

// NewNode constructor
func NewNode(address string) *Node {
	return &Node{
		routingTable: make(map[[28]byte]*Node),
		ID:           hashKey([]byte(address)),
		Addr:         address,
	}
}

// DHT distributed hash table
type DHT struct {
	id           [28]byte
	routingTable map[[28]byte]*Node
}

// NewDHT constructor
func NewDHT(address string) *DHT {

	id := hashKey([]byte(address))

	rt := make(map[[28]byte]*Node)
	// add self to routing table
	rt[id] = &Node{Addr: address, ID: id}

	return &DHT{
		id:           id,
		routingTable: rt,
	}
}

// Put returns which node to pass data too
func (d *DHT) Put(key []byte) (*Node, error) {

	keyHash := hashKey(key)

	closestNodes := d.findClosestNodes(keyHash, 0)
	if len(closestNodes) == 0 {
		return nil, &ErrKeyNotFound{Key: keyHash[:]}
	}

	var bestNode [28]byte

	for _, n := range closestNodes {
		if bestNode == [28]byte{} || compareDistances(xor(keyHash, n), xor(keyHash, bestNode)) < 0 {
			bestNode = n
		}
	}

	return d.routingTable[bestNode], nil
}

// AddNode to hash table
func (d *DHT) AddNode(n *Node) error {

	keyHash := hashKey(n.ID[:])
	n.ID = keyHash

	err := d.updateNodeRoutingTable(n, keyHash, d.routingTable)
	if err != nil {
		return fmt.Errorf("failed to update routing table %v", err)
	}

	return nil
}

// FindClosestNodes to first node
func (d *DHT) FindClosestNodes() [][28]byte {
	keyHash := hashKey(d.id[:])
	return d.findClosestNodes(keyHash, 1)[1:]
}

func (d *DHT) findClosestNodes(keyHash [28]byte, increment int) [][28]byte {

	var closest [][28]byte
	checked := make(map[[28]byte]struct{})

	// start with current DHT
	currentNode := d.id

	for i := 0; i < k+increment; i++ {

		var (
			bestDistance = [28]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
			bestNode = [28]byte{}
		)

		// iterate over all nodes relative to current node
		for i := range d.routingTable {
			if _, found := checked[i]; found {
				// skip if previously checked
				continue
			}

			distance := xor(keyHash, i)
			if compareDistances(distance, bestDistance) < 0 {
				bestDistance = distance
				bestNode = i
			}
		}

		closest = append(closest, bestNode)
		checked[currentNode] = struct{}{}
		currentNode = bestNode
	}

	return closest
}

func (d *DHT) updateNodeRoutingTable(node *Node, keyHash [28]byte, routingTable map[[28]byte]*Node) error {

	var closestNode *Node
	var closestDistance [28]byte

	for nodeID, existingNodeID := range d.routingTable {
		distance := xor(keyHash, nodeID)
		if closestNode == nil || compareDistances(distance, closestDistance) < 0 {
			closestNode = existingNodeID
			closestDistance = distance
		}
	}

	// set self in routing table
	node.routingTable = routingTable
	node.routingTable[keyHash] = node

	// TODO:
	// find nearest neighbors for node dissemination

	return nil
}

func hashKey(key []byte) [28]byte {
	h := sha3.New224()
	h.Write(key)
	hash := h.Sum(nil)

	var result [28]byte
	copy(result[:], hash[:28])

	return result
}

func xor(a, b [28]byte) [28]byte {
	var result [28]byte
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}

func compareDistances(a, b [28]byte) int {
	for i := range a {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
