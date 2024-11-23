package common

import (
	"context"
	"errors"
	"os"
	"sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	validAuthToken string
	initOnce       sync.Once // Ensures the token is initialized only once
)

// Initialize the validAuthToken
func initializeAuthToken() {
	validAuthToken = os.Getenv("AUTH_KEY")
	if validAuthToken == "" {
		panic("AUTH_KEY environment variable is not set")
	}
}

// UnaryInterceptor for Authorization
func AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Initialize the token only once
	initOnce.Do(initializeAuthToken)

	// Extract metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	// Check for authorization token
	authHeader, exists := md["authorization"]
	if !exists || len(authHeader) == 0 || authHeader[0] != validAuthToken {
		return nil, errors.New("unauthorized: invalid or missing authorization token")
	}

	// Proceed to the next handler if authorized
	return handler(ctx, req)
}
