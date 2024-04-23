package structs

import (
	"errors"
	"fmt"
	"net"
)

var (
	k = 3
)

// DHTNode
type DHTNode struct {
	ID   [28]byte
	Addr net.Addr
}

// DHT distributed hash table
type DHT struct {
	id           [28]byte
	addr         string // network address (IP:Port)
	routingTable map[[28]byte]*DHTNode
}

// NewDHT
func NewDHT(address string) (*DHT, error) {

	id := hashKey([]byte(address))
	fmt.Println(id)

	rt := make(map[[28]byte]*DHTNode)
	// add self to routing table
	rt[id] = &DHTNode{Addr: &net.IPAddr{}, ID: id}

	return &DHT{
		id:           id,
		routingTable: rt,
		addr:         address,
	}, nil
}

// Put
func (d *DHT) Put(key, data []byte) (*DHTNode, error) {

	keyHash := hashKey(key)

	closestNodes := d.findClosestNodes(keyHash, 0)
	if len(closestNodes) == 0 {
		return nil, errors.New("key not found")
	}

	var bestNode [28]byte

	for _, n := range closestNodes {
		if bestNode == [28]byte{} || compareDistances(xor(keyHash, n), xor(keyHash, bestNode)) < 0 {
			bestNode = n
		}
	}

	return d.routingTable[bestNode], nil
}

// AddNode
func (d *DHT) AddNode(n *DHTNode) {

}

// FindClosestNodes
func (d *DHT) FindClosestNodes() [][28]byte {
	keyHash := hashKey(d.id[:])
	return d.findClosestNodes(keyHash, 1)[1:]
}

func (d *DHT) findClosestNodes(keyHash [28]byte, increment int) [][28]byte {

	var closest [][28]byte
	checked := make(map[[28]byte]struct{})

	// start with current DHT
	currentNode := d.id
	fmt.Println(currentNode)
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
