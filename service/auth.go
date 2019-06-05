package service

import (
	"strings"

	"golang.org/x/net/context"

	"github.com/sirupsen/logrus"
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
	auth := md["authorization"]
	if len(auth) < 1 {
		return nil, ErrTokenInvalid
	}
	token := strings.TrimPrefix(auth[0], "Bearer ")
	logrus.Debug(token)
	// TODO: look for token in database and ensure it exists
	return nil, nil
}
