package rpcv1

import (
	"context"
	"math"
	"time"

	"github.com/ganitzsh/12fact/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func authorize(ctx context.Context, needRoot bool) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return service.ErrInvalidInput
	}
	authHeader, ok := md["authorization"]
	if !ok {
		return service.ErrTokenInvalid
	}
	return service.ValidateToken(authHeader[0])
}

func midValidateToken() grpc.ServerOption {
	return grpc.UnaryInterceptor(grpc.UnaryServerInterceptor(
		func(
			ctx context.Context,
			req interface{},
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler,
		) (interface{}, error) {
			needRoot := false
			if info.FullMethod == "/proto.IMage/GetNewToken" {
				needRoot = true
			}
			if err := authorize(ctx, needRoot); err != nil {
				return nil, err
			}
			return handler, nil
		}))
}

func midLogCall() grpc.ServerOption {
	return grpc.StreamInterceptor(grpc.StreamServerInterceptor(
		func(
			srv interface{},
			ss grpc.ServerStream,
			info *grpc.StreamServerInfo,
			handler grpc.StreamHandler,
		) error {
			md, _ := metadata.FromIncomingContext(ss.Context())
			start := time.Now()
			failed := false
			err := handler(srv, ss)
			end := time.Since(start)
			if err != nil {
				failed = true
			}
			latency := int(math.Ceil(float64(end.Nanoseconds()) / 1000000.0))
			logrus.WithFields(logrus.Fields{
				"path":       info.FullMethod,
				"start":      start,
				"end":        end,
				"latency":    latency,
				"user_agent": md["user-agent"],
				"failed":     failed,
			}).Info("New RPC server stream input")
			return err
		}))
}
