package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func validateToken(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrMissingGRPCAuthData
	}
	t := md["authorization"]
	if len(t) < 1 {
		return nil, ErrTokenInvalid
	}
	return nil, nil
}
