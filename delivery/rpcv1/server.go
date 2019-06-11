package rpcv1

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/ganitzsh/12fact/delivery/rpcv1/proto"
	"github.com/ganitzsh/12fact/service"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type RPCServer struct {
	*service.Config
}

func NewRPCServer(cfg *service.Config) *RPCServer {
	return &RPCServer{
		Config: cfg,
	}
}

func (s *RPCServer) ListenAndServe() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	if s.TLSEnabled {
		creds, err := credentials.NewServerTLSFromFile(s.TLSCert, s.TLSKey)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	opts = append(opts, midLogCall())
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterIMageServer(grpcServer, s)
	logrus.Info("RPC server started on port ", s.Port)
	return grpcServer.Serve(lis)
}

func (s *RPCServer) TransformImage(stream pb.IMage_TransformImageServer) error {
	var format string
	var any *any.Any
	var transformationType pb.TransformationType
	var fileSize int64

	f, e := ioutil.TempFile(os.TempDir(), "12fact")
	if e != nil {
		return e
	}
	if s.DevMode {
		logrus.WithField("file", f.Name()).Debug("Created temporary file")
		defer f.Close()
	} else {
		defer os.Remove(f.Name())
	}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Errorf("failed to receive stream: %v", err)
			return err
		}
		if req.GetImage().GetSize() > int64(s.MaxImageSize) {
			return service.ErrFileSizeExceeded
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
			logrus.Errorf("failed to write file: %v", err)
			return err
		}
	}
	format = strings.ToLower(format)
	logrus.WithFields(logrus.Fields{
		"type":      transformationType,
		"file_size": fileSize,
		"format":    format,
	}).Debug("Transforming image")
	if _, err := f.Seek(0, 0); err != nil {
		logrus.Errorf("failed to seek beigining of file: %v", err)
		return err
	}
	message := s.makeMessage(transformationType)
	if message == nil {
		return service.ErrInternalError
	}
	if err := ptypes.UnmarshalAny(any, message); err != nil {
		return err
	}
	ret, err := service.SingleTransformImage(
		f, format, s.ToTransfomation(message),
	)
	if err != nil {
		logrus.Errorf("failed to transform image: %v", err)
		return err
	}
	if ret == nil {
		logrus.Errorf("Something went wrong")
		return service.ErrInternalError
	}
	for {
		buff := make([]byte, s.BufferSize)
		_, err := ret.Read(buff)
		if err == io.EOF {
			break
		} else if err != nil {
			logrus.Errorf("failed to read image buffer: %v", err)
			return err
		}
		if e := stream.Send(&pb.TransformedImage{
			File: buff,
		}); e != nil {
			logrus.Errorf("failed to send file data: %v", e)
			return e
		}
	}
	return nil
}
