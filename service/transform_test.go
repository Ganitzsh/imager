package service_test

import (
	"os"
	"testing"

	"github.com/ganitzsh/12fact/service"
	"github.com/stretchr/testify/require"
)

var testFiles = []string{"../example/img.png"}

func TestTransformNils(t *testing.T) {
	trans1 := service.NewRotate().SetAngle(90).SetClockWise(true)

	f, err := os.Open(testFiles[0])
	defer f.Close()
	require.NoError(t, err)

	ret, err := service.SingleTransformImage(f, ".png", nil)
	require.Error(t, err)
	require.Nil(t, ret)

	ret, err = service.SingleTransformImage(nil, "", nil)
	require.Error(t, err)
	require.Nil(t, ret)

	ret, err = service.SingleTransformImage(nil, ".tmp", nil)
	require.Error(t, err)
	require.Nil(t, ret)

	ret, err = service.TransformImage(f, ".nope", []service.Transformation{
		trans1,
	})
	require.Error(t, err)
	require.Nil(t, ret)
}

type testTransform struct {
	trans      service.Transformation
	shouldFail bool
}

func TestRotate(t *testing.T) {
	transformations := []testTransform{
		testTransform{
			trans: service.NewRotate().SetAngle(45).SetClockWise(false),
		}, testTransform{
			trans: service.NewRotate().SetAngle(90).SetClockWise(true),
		}, testTransform{
			trans: service.NewRotate().SetAngle(-45),
		}, testTransform{
			trans: service.NewRotate().SetAngle(1.34),
		}, testTransform{
			trans: service.NewRotate().SetAngle(1234.890),
		}, testTransform{
			trans: service.NewRotate().SetAngle(90).SetClockWise(true),
		},
	}

	f, err := os.Open(testFiles[0])
	defer f.Close()
	require.NoError(t, err)

	for _, tc := range transformations {
		f.Seek(0, 0)
		ret, err := service.SingleTransformImage(f, ".png", tc.trans)
		if tc.shouldFail {
			require.Error(t, err)
			require.Nil(t, ret)
		} else {
			require.NoError(t, err)
			require.NotNil(t, ret)
		}
	}
}

func TestBlur(t *testing.T) {
	transformations := []testTransform{
		testTransform{
			trans: service.NewBlur().SetSigma(3.4),
		}, testTransform{
			trans: service.NewBlur().SetSigma(9.0),
		},
		// testTransform{
		// 	trans: service.NewBlur().SetSigma(9999.9999),
		// },
		testTransform{
			trans: service.NewBlur().SetSigma(-45),
		},
	}

	f, err := os.Open(testFiles[0])
	defer f.Close()
	require.NoError(t, err)

	for _, tc := range transformations {
		f.Seek(0, 0)
		ret, err := service.SingleTransformImage(f, ".png", tc.trans)
		if tc.shouldFail {
			require.Error(t, err)
			require.Nil(t, ret)
		} else {
			require.NoError(t, err)
			require.NotNil(t, ret)
		}
	}
}

func TestCrop(t *testing.T) {
	transformations := []testTransform{
		testTransform{
			trans: service.NewCrop().
				SetWidth(123).SetHeight(123).SetTopLeftX(0).SetTopLeftY(0),
		}, testTransform{
			trans: service.NewCrop().
				SetWidth(0).SetHeight(123).SetTopLeftX(0).SetTopLeftY(0),
			shouldFail: true,
		}, testTransform{
			trans: service.NewCrop().
				SetWidth(50000).SetHeight(345).SetTopLeftX(0).SetTopLeftY(0),
		}, testTransform{
			trans: service.NewCrop().
				SetWidth(0).SetHeight(0).SetTopLeftX(0).SetTopLeftY(0),
			shouldFail: true,
		},
	}

	for _, fp := range testFiles {
		f, err := os.Open(fp)
		require.NoError(t, err)
		for _, tc := range transformations {
			f.Seek(0, 0)
			ret, err := service.SingleTransformImage(f, ".png", tc.trans)
			if tc.shouldFail {
				require.Error(t, err)
				require.Nil(t, ret)
			} else {
				require.NoError(t, err)
				require.NotNil(t, ret)
			}
		}
		f.Close()
	}
}

func TestResize(t *testing.T) {
	transformations := []testTransform{
		testTransform{
			trans: service.NewResize().SetWidth(540).SetHeight(1090),
		}, testTransform{
			trans: service.NewResize().SetWidth(678).SetHeight(0),
		}, testTransform{
			trans: service.NewResize().SetWidth(0).SetHeight(678),
		}, testTransform{
			trans: service.NewResize().SetWidth(678).SetHeight(0),
		}, testTransform{
			trans:      service.NewResize().SetWidth(0).SetHeight(0),
			shouldFail: true,
		}, testTransform{
			trans:      service.NewResize().SetWidth(-123).SetHeight(490),
			shouldFail: true,
		}, testTransform{
			trans:      service.NewResize().SetWidth(0).SetHeight(-490),
			shouldFail: true,
		}, testTransform{
			trans:      service.NewResize().SetWidth(-5000).SetHeight(-10),
			shouldFail: true,
		},
	}

	f, err := os.Open(testFiles[0])
	defer f.Close()
	require.NoError(t, err)

	for _, tc := range transformations {
		f.Seek(0, 0)
		ret, err := service.SingleTransformImage(f, ".png", tc.trans)
		if tc.shouldFail {
			require.Error(t, err)
			require.Nil(t, ret)
		} else {
			require.NoError(t, err)
			require.NotNil(t, ret)
		}
	}
}
