package service_test

import (
	"os"
	"testing"

	"github.com/ganitzsh/12fact/service"
	"github.com/stretchr/testify/require"
)

func TestRotate(t *testing.T) {
	trans1 := service.NewRotate().SetAngle(45).SetClockWise(false)
	trans2 := service.NewRotate().SetAngle(90).SetClockWise(true)

	f, err := os.Open("../example/img.png")
	defer f.Close()
	require.NoError(t, err)
	ret, err := service.SingleTransformImage(f, ".png", trans1)
	require.NoError(t, err)
	require.NotNil(t, ret)

	f.Seek(0, 0)
	ret, err = service.SingleTransformImage(f, ".png", nil)
	require.Error(t, err)
	require.Nil(t, ret)

	ret, err = service.SingleTransformImage(nil, "", nil)
	require.Error(t, err)
	require.Nil(t, ret)

	ret, err = service.SingleTransformImage(nil, ".tmp", nil)
	require.Error(t, err)
	require.Nil(t, ret)

	ret, err = service.SingleTransformImage(f, ".png", trans2)
	require.NoError(t, err)
	require.NotNil(t, ret)

	f.Seek(0, 0)
	ret, err = service.TransformImage(f, ".png", []service.Transformation{
		trans1, trans2,
	})
	require.NoError(t, err)
	require.NotNil(t, ret)

	ret, err = service.TransformImage(f, ".nope", []service.Transformation{
		trans1, trans2,
	})
	require.Error(t, err)
	require.Nil(t, ret)
}

func TestBlur(t *testing.T) {
	trans1 := service.NewBlur().SetSigma(3.4)
	trans2 := service.NewBlur().SetSigma(9.0)

	f, err := os.Open("../example/img.png")
	defer f.Close()
	require.NoError(t, err)
	ret, err := service.SingleTransformImage(f, ".png", trans1)
	require.NoError(t, err)
	require.NotNil(t, ret)

	f.Seek(0, 0)
	ret, err = service.SingleTransformImage(f, ".png", nil)
	require.Error(t, err)
	require.Nil(t, ret)

	ret, err = service.SingleTransformImage(nil, "", nil)
	require.Error(t, err)
	require.Nil(t, ret)

	ret, err = service.SingleTransformImage(nil, ".tmp", nil)
	require.Error(t, err)
	require.Nil(t, ret)

	f.Seek(0, 0)
	ret, err = service.SingleTransformImage(f, ".png", trans2)
	require.NoError(t, err)
	require.NotNil(t, ret)

	f.Seek(0, 0)
	ret, err = service.TransformImage(f, ".png", []service.Transformation{
		trans1, trans2,
	})
	require.NoError(t, err)
	require.NotNil(t, ret)

	ret, err = service.TransformImage(f, ".nope", []service.Transformation{
		trans1, trans2,
	})
	require.Error(t, err)
	require.Nil(t, ret)
}

func TestCrop(t *testing.T) {
}
