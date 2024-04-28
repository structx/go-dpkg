package domain

import (
	"context"

	"google.golang.org/grpc"
)

// AuthInterceptor gRPC auth interceptor interface
//
//go:generate mockery --name AuthInterceptor
type AuthInterceptor interface {
	// UnaryInterceptor single request wallet verification
	UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	// StreamInterceptor streaming wallet verification
	StreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
}
