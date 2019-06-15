package rpcv1

import (
	"sync"

	pb "github.com/ganitzsh/12fact/delivery/rpcv1/proto"
	"github.com/ganitzsh/12fact/service"
	"github.com/golang/protobuf/proto"
)

var mutDecoder sync.Mutex
var requestFactory = map[pb.TransformationType]func() proto.Message{
	pb.TransformationType_ROTATE: func() proto.Message {
		return &pb.RotateImageRequest{}
	},
	pb.TransformationType_BLUR: func() proto.Message {
		return &pb.BlurImageRequest{}
	},
	pb.TransformationType_CROP: func() proto.Message {
		return &pb.CropImageRequest{}
	},
	pb.TransformationType_RESIZE: func() proto.Message {
		return &pb.ResizeImageRequest{}
	},
}

func (s *RPCServer) makeMessage(typ pb.TransformationType) proto.Message {
	mutDecoder.Lock()
	fn := requestFactory[typ]
	if fn == nil {
		return nil
	}
	ret := fn()
	mutDecoder.Unlock()
	return ret
}

func (s *RPCServer) ToTransfomation(req proto.Message) service.Transformation {
	if req == nil {
		return nil
	}
	switch req.(type) {
	case *pb.RotateImageRequest:
		return service.NewRotate().
			SetAngle(float64(req.(*pb.RotateImageRequest).GetAngle())).
			SetClockWise(req.(*pb.RotateImageRequest).GetClockWise())
	case *pb.BlurImageRequest:
		return service.NewBlur().
			SetSigma(float64(req.(*pb.BlurImageRequest).GetSigma()))
	case *pb.CropImageRequest:
		return service.NewCrop().
			SetTopLeftX(int(req.(*pb.CropImageRequest).GetTopLeftX())).
			SetTopLeftY(int(req.(*pb.CropImageRequest).GetTopLeftY())).
			SetWidth(int(req.(*pb.CropImageRequest).GetWidth())).
			SetHeight(int(req.(*pb.CropImageRequest).GetHeight()))
	case *pb.ResizeImageRequest:
		return service.NewResize().
			SetWidth(int(req.(*pb.ResizeImageRequest).GetWidth())).
			SetHeight(int(req.(*pb.ResizeImageRequest).GetHeight()))
	default:
		return nil
	}
}
