package rpcv1_test

import (
	"testing"

	"github.com/ganitzsh/12fact/delivery/rpcv1"
	pb "github.com/ganitzsh/12fact/delivery/rpcv1/proto"
	"github.com/ganitzsh/12fact/service"
	"github.com/stretchr/testify/require"
)

func TestToTransformation(t *testing.T) {
	s := rpcv1.NewRPCServer(nil)
	rotate := &pb.RotateImageRequest{
		Angle:     45,
		ClockWise: true,
	}
	transformedRotate := s.ToTransfomation(rotate)
	require.NotNil(t, transformedRotate)
	require.EqualValues(t, &service.Rotate{
		Angle:     45,
		ClockWise: true,
	}, transformedRotate)
	blur := &pb.BlurImageRequest{
		Sigma: 2.0,
	}
	transformedBlur := s.ToTransfomation(blur)
	require.NotNil(t, transformedBlur)
	require.EqualValues(t, &service.Blur{
		Sigma: 2.0,
	}, transformedBlur)
	crop := &pb.CropImageRequest{
		TopLeftX: 590,
		TopLeftY: 104,
		Width:    450,
		Height:   900,
	}
	transformedCrop := s.ToTransfomation(crop)
	require.NotNil(t, transformedCrop)
	require.EqualValues(t, &service.Crop{
		TopLeftX: 590,
		TopLeftY: 104,
		Width:    450,
		Height:   900,
	}, transformedCrop)
	resize := &pb.ResizeImageRequest{
		Width:  450,
		Height: 900,
	}
	transformedResize := s.ToTransfomation(resize)
	require.NotNil(t, transformedResize)
	require.EqualValues(t, &service.Resize{
		Width:  450,
		Height: 900,
	}, transformedResize)
}
