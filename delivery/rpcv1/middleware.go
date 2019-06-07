package rpcv1

import (
	"math"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

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
