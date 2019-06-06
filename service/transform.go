package service

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"
	"os"

	"github.com/disintegration/imaging"
	pb "github.com/ganitzsh/12fact/proto"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/sirupsen/logrus"
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

func GetFormatFromExtension(ext string) imaging.Format {
	return formats[ext]
}

func supportedExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".tif", ".tiff", ".bmp", ".gif":
		return true
	default:
		return false
	}
}

func TransformImageFunc(
	f *os.File, ext string,
	typ TransformationType,
	data *any.Any,
) (io.Reader, error) {
	if !supportedExt(ext) {
		return nil, ErrUnsupportedFormat
	}
	img, err := imaging.Decode(f)
	if err != nil {
		logrus.Errorf("failed to decode image: %v", err)
		return nil, err
	}
	var fn func(image.Image, imaging.Format, proto.Message) (io.Reader, error)
	var out proto.Message
	switch typ {
	case pTransformationTypeRotate:
		out = &pb.RotateImageRequest{}
		fn = Rotate
		break
	case pb.TransformationType_CROP:
		out = &pb.CropImageRequest{}
		fn = Crop
		break
	case pb.TransformationType_BLUR:
		out = &pb.BlurImageRequest{}
		fn = Blur
		break
	default:
		return nil, errors.New("Unkown transformation type")
	}
	if err := ptypes.UnmarshalAny(data, out); err != nil {
		return nil, err
	}
	return fn(img, formats[ext], out)
}

func Rotate(
	img image.Image,
	format imaging.Format,
	req interface{},
) (io.Reader, error) {
	direction := -1.0
	args := req.(*pb.RotateImageRequest)
	if args.GetClockWise() {
		direction *= 1.0
	}
	return imageBuffer(imaging.Rotate(
		img, float64(args.GetAngle())*direction, color.Black,
	), format)
}

type CropOptions struct {
	TopLeftX int
	TopLeftY int
	Width    int
	Height   int
}

func Crop(
	img image.Image,
	format imaging.Format,
	req interface{},
) (io.Reader, error) {
	args, ok := req.(*CropOptions)
	if !ok {
		return nil, ErrInternalError
	}
	return imageBuffer(imaging.Crop(img, image.Rect(
		args.TopLeftX, args.TopLeftY,
		args.TopLeftX+args.Width, args.TopLeftY+args.Height,
	)), format)
}

type CropOptions struct {
	TopLeftX int
}

func Blur(
	img image.Image,
	format imaging.Format,
	req interface{},
) (io.Reader, error) {
	return imageBuffer(imaging.Blur(
		img, float64(req.(*pb.BlurImageRequest).GetSigma()),
	), format)
}

func imageBuffer(img image.Image, format imaging.Format) (io.Reader, error) {
	ret := bytes.Buffer{}
	if err := imaging.Encode(&ret, img, format); err != nil {
		return nil, err
	}
	return &ret, nil
}
