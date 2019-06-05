package service

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/disintegration/imaging"
	pb "github.com/ganitzsh/12fact/proto"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

var formats = map[string]imaging.Format{
	".jpg":  imaging.JPEG,
	".jpeg": imaging.JPEG,
	".png":  imaging.PNG,
	".tif":  imaging.TIFF,
	".tiff": imaging.TIFF,
	".bmp":  imaging.BMP,
	".gif":  imaging.GIF,
}

func (s *Server) transformImage(
	f *os.File, ext string,
	typ pb.TransformationType,
	data *any.Any,
) (io.Reader, error) {
	if formats[ext] == 0 {
		return nil, ErrUnsupportedFormat
	}
	img, err := imaging.Decode(f)
	if err != nil {
		return nil, err
	}
	var fn func(image.Image, imaging.Format, interface{}) (io.Reader, error)
	var out proto.Message
	switch typ {
	case pb.TransformationType_ROTATE:
		out = &pb.RotateImageRequest{}
		fn = rotate
		break
	case pb.TransformationType_CROP:
		out = &pb.CropImageRequest{}
		fn = crop
		break
	case pb.TransformationType_BLUR:
		out = &pb.BlurImageRequest{}
		fn = blur
		break
	default:
		return nil, errors.New("Unkown transformation type")
	}
	if err := ptypes.UnmarshalAny(data, out); err != nil {
		return nil, err
	}
	if s.DevMode {
		spew.Dump(out)
	}
	return fn(img, formats[ext], &out)
}

func rotate(
	img image.Image,
	format imaging.Format,
	req interface{},
) (io.Reader, error) {
	return imageBuffer(imaging.Rotate(
		img, float64(req.(*pb.RotateImageRequest).GetAngle()), color.Black,
	), format)
}

func crop(
	img image.Image,
	format imaging.Format,
	req interface{},
) (io.Reader, error) {
	return nil, nil
}

func blur(
	img image.Image,
	format imaging.Format,
	req interface{},
) (io.Reader, error) {
	return nil, nil
}

func imageBuffer(img image.Image, format imaging.Format) (io.Reader, error) {
	ret := bytes.Buffer{}
	if err := imaging.Encode(&ret, img, format); err != nil {
		return nil, err
	}
	return &ret, nil
}
