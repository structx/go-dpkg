// Package messagebroker implementation
package messagebroker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/structx/go-dpkg/domain"
	pbv1 "github.com/structx/go-dpkg/proto/messaging/v1"
)

// Client message broker implementation
type Client struct {
	conn *grpc.ClientConn
}

// interface compliance
var _ domain.MessageBroker = (*Client)(nil)

// New constructor
func New(cfg domain.Config) (*Client, error) {

	mcfg := cfg.GetMessenger()
	if mcfg == nil {
		return nil, errors.New("missing message broker configuration")
	}

	conn, err := grpc.Dial(mcfg.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial message broker %v", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

// Publish message to topic
func (c *Client) Publish(ctx context.Context, topic string, msg []byte) error {
	c.conn.Connect()

	cli := pbv1.NewMessagingServiceV1Client(c.conn)

	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	_, err := cli.Publish(timeout, &pbv1.Envelope{
		Topic:   topic,
		Payload: msg,
	})
	if err != nil {
		return fmt.Errorf("failed to publish message to %s %v", topic, err)
	}

	return nil
}

// Subscribe to topic
func (c *Client) Subscribe(ctx context.Context, topic string) (<-chan domain.Envelope, error) {
	c.conn.Connect()

	cli := pbv1.NewMessagingServiceV1Client(c.conn)

	stream, err := cli.Subscribe(ctx, &pbv1.Subscription{
		Topic: topic,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to subscribe to %s %v", topic, err)
	}

	ch := make(chan domain.Envelope)
	errCh := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(errCh)

		for {

			select {
			case <-ctx.Done():
				err := stream.CloseSend()
				if err != nil {
					errCh <- fmt.Errorf("failed to close send stream %v", err)
				}
				return
			default:

				m, err := stream.Recv()
				if err != nil {
					errCh <- fmt.Errorf("failed to receive message %v", err)
					return
				}

				ch <- domain.NewMsg(
					m.GetTopic(),
					m.GetPayload(),
				)

			}
		}

	}()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				return nil, fmt.Errorf("failed to keep subscription alive %v", err)
			}
		default:
			return ch, nil
		}
	}
}

// RequestResponse send request and wait for response
func (c *Client) RequestResponse(ctx context.Context, in domain.Envelope) (domain.Envelope, error) {
	c.conn.Connect()

	cli := pbv1.NewMessagingServiceV1Client(c.conn)

	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	rr, err := cli.RequestResponse(timeout, &pbv1.Envelope{
		Topic:   in.GetTopic(),
		Payload: in.GetPayload(),
	})
	if err != nil {
		return nil, fmt.Errorf("response/response message failed %v", err)
	}

	return domain.NewMsg(rr.GetTopic(), rr.GetPayload()), nil
}

// Close message broker connections
func (c *Client) Close() error {
	return c.conn.Close()
}
