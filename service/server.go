package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"github.com/davecgh/go-spew/spew"
	pb "github.com/ganitzsh/12fact/proto"
	"github.com/golang/protobuf/ptypes/any"
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

func (s *Server) TransformImage(stream pb.IMage_TransformImageServer) error {
	var format string
	var any *any.Any
	var transformationType pb.TransformationType
	var fileSize int64

	f, e := ioutil.TempFile(os.TempDir(), "12fact")
	if e != nil {
		return e
	}
	logrus.Debug(f.Name())
	defer f.Close()
	// defer os.Remove(f.Name())
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			logrus.Error("1 ", err)
			return err
		}
		if req.GetImage().GetSize() > s.MaxImageSize {
			return ErrFileSizeExceeded
		}
		if fileSize == 0 {
			fileSize = req.GetImage().GetSize()
		}
		if format == "" {
			format = req.GetImage().GetFormat()
		}
		if any == nil {
			any = req.GetData()
		}
		if transformationType == 0 {
			transformationType = req.GetType()
		}
		_, err = f.Write(req.GetImage().File)
		if err != nil {
			logrus.Error("2 ", err)
			return err
		}
		// logrus.Infof("Writing chunk of %d bytes", n)
	}
	logrus.WithFields(logrus.Fields{
		"type":      transformationType,
		"data":      any,
		"file_size": fileSize,
	}).Debug("Transforming image")
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	if s.DevMode {
		spew.Dump(stat)
	}
	return nil
}
