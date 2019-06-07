package service

import (
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
)

type Rotate struct {
	Transformation

	Angle     float64
	ClockWise bool
}

func NewRotate() *Rotate {
	return &Rotate{}
}

func (t *Rotate) SetAngle(value float64) *Rotate {
	t.Angle = value
	return t
}

func (t *Rotate) SetClockWise(value bool) *Rotate {
	t.ClockWise = value
	return t
}

func (Rotate) GetType() TransformationType {
	return TransformationTypeRotate
}

func (t *Rotate) Log() {
	LogTransformation(t, logrus.Fields{
		"angle":      t.Angle,
		"clock_wise": t.ClockWise,
	})
}

func (t *Rotate) Do(img *Image) *Image {
	direction := 1.0
	if t.ClockWise {
		direction = -1.0
	}
	finalAngle := t.Angle * direction
	img.Image = imaging.Rotate(img, finalAngle, color.Black)
	return img
}
