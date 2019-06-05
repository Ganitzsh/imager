package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"github.com/davecgh/go-spew/spew"
	pb "github.com/ganitzsh/12fact/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	*Config
}

func NewServer() (*Server, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}
	return &Server{
		Config: cfg,
	}, nil
}

func (s *Server) ListenAndServe() error {
	if s.DevMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.Port))
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	if s.TLSEnabled {
		// if *certFile == "" {
		// 	*certFile = testdata.Path("server1.pem")
		// }
		// if *keyFile == "" {
		// 	*keyFile = testdata.Path("server1.key")
		// }
		// creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		// if err != nil {
		// 	log.Fatalf("Failed to generate credentials %v", err)
		// }
		// opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterIMageServer(grpcServer, s)
	logrus.Info("Listening on port ", s.Port)
	return grpcServer.Serve(lis)
}

func (s *Server) RotateImage(stream pb.IMage_RotateImageServer) error {
	var format string
	var angle int32

	f, e := ioutil.TempFile(os.TempDir(), "12fact")
	if e != nil {
		return e
	}
	logrus.Debug(f.Name())
	defer os.Remove(f.Name())
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if format == "" {
			format = req.GetImage().GetFormat()
		}
		if angle == 0 {
			angle = req.GetAngle()
		}
		n, err := f.Write(req.GetImage().File)
		if err != nil {
			return err
		}
		fmt.Printf("Reveived %d bytes from client\n", n)
	}
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	spew.Dump(stat)
	return nil
}

func (s *Server) BlurImage(stream pb.IMage_BlurImageServer) error {
	return nil
}

func (s *Server) CropImage(stream pb.IMage_CropImageServer) error {
	return nil
}
