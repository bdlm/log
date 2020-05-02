// Package interceptor provides gRPC interceptors for common middleware
// requirements.
package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

// Middleware defines the interface for gRPC interceptor handlers.
type Middleware interface {
	StreamInterceptor(
		interface{},
		grpc.ServerStream,
		*grpc.StreamServerInfo,
		grpc.StreamHandler,
	) error

	UnaryInterceptor(
		context.Context,
		interface{},
		*grpc.UnaryServerInfo,
		grpc.UnaryHandler,
	) (interface{}, error)
}
