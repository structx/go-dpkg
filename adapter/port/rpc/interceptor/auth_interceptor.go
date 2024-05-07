// Package interceptor gRPC request interceptors
package interceptor

import (
	"context"
	"encoding/json"

	"github.com/structx/go-dpkg/adapter/port/http/header"
	"github.com/structx/go-dpkg/domain"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type wrappedStream struct {
	grpc.ServerStream
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

// RecvMsg ...
func (w *wrappedStream) RecvMsg(m any) error {
	return w.ServerStream.RecvMsg(m)
}

// SendMessage ...
func (w *wrappedStream) SendMsg(m any) error {
	return w.ServerStream.SendMsg(m)
}

// Auth interceptor implementation
type Auth struct {
	log *zap.SugaredLogger
	mb  domain.MessageBroker
}

// NewAuth constructor
func NewAuth(logger *zap.Logger, broker domain.MessageBroker) *Auth {
	return &Auth{
		log: logger.Sugar().Named("AuthInterceptor"),
		mb:  broker,
	}
}

// UnaryInterceptor single request interceptor to verify wallet access permissions
func (a *Auth) UnaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrMissingMetadata
	}

	v := valid(md[header.DecentralizedIdentity])
	if !v {
		return nil, ErrInvalidToken
	}

	msg := domain.NewMsg(domain.VerifyServiceAccess.String(), []byte{})
	rr, err := a.mb.RequestResponse(ctx, msg)
	if err != nil {
		a.log.Errorf("request response failed %v", err)
		return nil, status.Error(codes.Internal, "access control check failed")
	}

	var en domain.EntryResponse
	err = json.Unmarshal(rr.GetPayload(), &en)
	if err != nil {
		a.log.Errorf("failed to unmarshal entry response %v", err)
		return nil, status.Error(codes.Internal, "access control check failed")
	}

	if !en.Granted {
		return nil, status.Error(codes.PermissionDenied, "access denied")
	}

	m, err := handler(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "RPC failed with error %v", err)
	}

	return m, nil
}

// StreamInterceptor streaming request interceptor to verify wallet access permissions
func (a *Auth) StreamInterceptor(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	ctx := ss.Context()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ErrMissingMetadata
	}

	v := valid(md[header.DecentralizedIdentity])
	if !v {
		return ErrInvalidToken
	}

	msg := domain.NewMsg(domain.VerifyServiceAccess.String(), []byte{})
	rr, err := a.mb.RequestResponse(ctx, msg)
	if err != nil {
		a.log.Errorf("request response failed %v", err)
		return status.Error(codes.Internal, "access control check failed")
	}

	var en domain.EntryResponse
	err = json.Unmarshal(rr.GetPayload(), &en)
	if err != nil {
		a.log.Errorf("failed to unmarshal entry response %v", err)
		return status.Error(codes.Internal, "access control check failed")
	}

	if !en.Granted {
		return status.Error(codes.PermissionDenied, "access denied")
	}

	err = handler(srv, newWrappedStream(ss))
	if err != nil {
		return status.Errorf(codes.Internal, "RPC failed with error %v", err)
	}

	return nil
}

func valid(s []string) bool {
	return len(s) < 1
}
