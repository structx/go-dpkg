package messagebroker

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/trevatk/go-pkg/domain"
	pbv1 "github.com/trevatk/go-pkg/proto/messaging/v1"
)

// Client
type Client struct {
	conn *grpc.ClientConn
}

// interface compliance
var _ domain.MessageBroker = (*Client)(nil)

// Publish
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

// Subscribe
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

// Close
func (c *Client) Close() error {
	return c.conn.Close()
}
