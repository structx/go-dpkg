// Package dht gRPC implementation
package dht

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/structx/go-pkg/domain"
	pbv1 "github.com/structx/go-pkg/proto/dht/v1"
)

// Client implementation
type Client struct {
	conn *grpc.ClientConn
}

// NewClient constructor
func NewClient(address string) (*Client, error) {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial client %v", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

// FindNode gRPC client call
func (c *Client) FindNode(ctx context.Context, nodeID, sender domain.NodeID) ([]*domain.Contact, error) {
	c.conn.Connect()

	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	cli := pbv1.NewDHTServiceClient(c.conn)
	response, err := cli.FindNode(timeout, &pbv1.FindNodeRequest{
		Sender: &pbv1.Sender{
			SenderId:    sender[:],
			RequestedAt: timestamppb.Now(),
		},
		NodeId: nodeID[:],
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send find node request %v", err)
	}

	contactSlice := make([]*domain.Contact, 0)
	for _, c := range response.ContactList {
		contactSlice = append(contactSlice, &domain.Contact{
			IP:   c.Ip,
			Port: int(c.GetPort()),
			ID:   domain.NodeID(c.NodeId),
		})
	}

	return contactSlice, nil
}

// Close client connection
func (c *Client) Close() error {
	return c.conn.Close()
}
