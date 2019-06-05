package service

import (
	pb "github.com/ganitzsh/12fact/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	pb.IMageClient

	ServerAddr string
	TLSCert    string
}

func NewClient() (*Client, error) {
	serverAddr := viper.GetString("Host")
	tlsCert := viper.GetString("TLSCert")
	var opts []grpc.DialOption
	if tlsCert != "" {
		creds, err := credentials.NewClientTLSFromFile(tlsCert, "")
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		ServerAddr:  serverAddr,
		TLSCert:     tlsCert,
		IMageClient: pb.NewIMageClient(conn),
	}, nil
}
