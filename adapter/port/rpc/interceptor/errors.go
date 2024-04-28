package interceptor

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrInvalidToken missing token
	ErrInvalidToken = status.Error(codes.InvalidArgument, "invalid token")
	// ErrMissingMetadata invalid or missing metadata
	ErrMissingMetadata = status.Error(codes.InvalidArgument, "missing metadata")
)
