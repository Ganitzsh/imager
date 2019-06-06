package service

import (
	"github.com/disintegration/imaging"
)

type Blur struct {
	Transformation

	Sigma float64
}

func NewBlur() *Blur {
	return &Blur{}
}

func (t *Blur) SetSigma(value float64) *Blur {
	t.Sigma = value
	return t
}

func (Blur) GetType() TransformationType {
	return TransformationTypeBlur
}

func (t *Blur) Do(img *Image) *Image {
	img.Image = imaging.Blur(img, t.Sigma)
	return img
}
