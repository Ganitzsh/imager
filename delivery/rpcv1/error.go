package rpcv1

import "github.com/ganitzsh/12fact/service"

var (
	ErrMissingGRPCAuthData = service.NewServiceError(
		"grpc_missing_data",
		"Missing auth data in GRPC call",
	)
)
