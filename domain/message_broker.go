package domain

import (
	"context"
)

// Topic pub/sub topics
type Topic string

// ServiceType supported service types
type ServiceType string

const (
	// JoinedRaftV1 joined network v1 topic
	JoinedRaftV1 Topic = "msg.v1.joined_network"

	// Mora message broker service
	Mora ServiceType = "mora"
)

// String stringify topic
func (t Topic) String() string {
	switch t {
	case JoinedRaftV1:
		return string(JoinedRaftV1)
	default:
		return ""
	}
}

// MessageBroker interface
type MessageBroker interface {
	// Publish message to topic
	Publish(context.Context, string, []byte) error
	// Subscribe to topic
	Subscribe(context.Context, string) (<-chan Envelope, error)
	// Close message broker connection
	Close() error
}

// Envelope msg interface
type Envelope interface {
	// GetTopic getter topic
	GetTopic() string
	// GetPayload getter payload
	GetPayload() []byte
}

// Msg envelope implementation
type Msg struct {
	topic   string
	payload []byte
}

// interface compliance
var _ Envelope = (*Msg)(nil)

// NewMsg constructor
func NewMsg(topic string, payload []byte) *Msg {
	return &Msg{
		topic:   topic,
		payload: payload,
	}
}

// GetTopic getter topic
func (m *Msg) GetTopic() string {
	return m.topic
}

// GetPayload getter payload
func (m *Msg) GetPayload() []byte {
	return m.payload
}

// JoinedNetwork new service joined network message
type JoinedNetwork struct {
	ServiceType ServiceType `json:"service_type"`
	ServiceID   string      `json:"service_id"`
}

// NewJoinedNetwork constructor
func NewJoinedNetwork(serviceID string, serviceType ServiceType) *JoinedNetwork {
	return &JoinedNetwork{
		ServiceType: serviceType,
		ServiceID:   serviceID,
	}
}
