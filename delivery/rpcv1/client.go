package rpcv1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	pb "github.com/ganitzsh/12fact/delivery/rpcv1/proto"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	pb.IMageClient

	ServerAddr string
	TLSCert    string

	bufferSize uint32
}

func NewClient() (*Client, error) {
	serverAddr := fmt.Sprintf(
		"%s:%d", viper.GetString("host"), viper.GetInt("port"),
	)
	tlsCert := viper.GetString("tls.cert")
	tlsEnabled := viper.GetBool("tls.enabled")
	var opts []grpc.DialOption
	if tlsEnabled {
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
	bufferSize := viper.GetUint32("buffer_size")
	if bufferSize == 0 {
		logrus.Warn("Invalid buffer size, setting to 2048 bytes")
		bufferSize = 2048
	}
	return &Client{
		ServerAddr:  serverAddr,
		TLSCert:     tlsCert,
		IMageClient: pb.NewIMageClient(conn),

		bufferSize: bufferSize,
	}, nil
}

func openFile(file string) (f *os.File, ext string, size int64, err error) {
	ext = filepath.Ext(file)
	if ext == "" {
		err = errors.New("Unknown extension")
		return
	}
	f, err = os.Open(file)
	if err != nil {
		return
	}
	stat, err := f.Stat()
	if err != nil {
		return
	}
	size = stat.Size()
	return
}

func (client *Client) streamTransform(
	f io.Reader,
	stream pb.IMage_TransformImageClient,
	req *pb.TransformImageRequest,
) (io.Reader, error) {
	if e := stream.Send(req); e != nil {
		return nil, e
	}
	buff := make([]byte, client.bufferSize)
	for {
		_, e := f.Read(buff)
		if e == io.EOF {
			break
		}
		if e := stream.Send(&pb.TransformImageRequest{
			Image: &pb.Image{
				File: buff,
			},
		}); e != nil {
			return nil, e
		}
	}
	if e := stream.CloseSend(); e != nil {
		return nil, e
	}

	outFile := new(bytes.Buffer)
	for {
		out, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if _, err := outFile.Write(out.GetFile()); err != nil {
			return nil, err
		}
	}
	return outFile, nil
}

func (client *Client) Transform(
	file string,
	typ pb.TransformationType,
	req proto.Message,
) (io.Reader, error) {
	f, ext, size, err := openFile(file)
	if err != nil {
		return nil, err
	}
	stream, err := client.TransformImage(context.Background())
	if err != nil {
		return nil, err
	}
	data, err := ptypes.MarshalAny(req)
	if err != nil {
		return nil, err
	}
	return client.streamTransform(f, stream, &pb.TransformImageRequest{
		Type: typ,
		Data: data,
		Image: &pb.Image{
			Size:   size,
			Format: ext,
		},
	})
}

func (c *Client) SaveToFile(r io.Reader, path string) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"file":  path,
		"wrote": fmt.Sprintf("%d bytes", n),
	}).Info("Writing file to disk")
	return ioutil.WriteFile(path, data, os.ModeAppend)
}
