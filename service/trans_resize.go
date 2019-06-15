package service

import (
	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
)

type Resize struct {
	Transformation

	Width  int
	Height int
}

func NewResize() *Resize {
	return &Resize{}
}

func (t *Resize) SetWidth(value int) *Resize {
	t.Width = value
	return t
}

func (t *Resize) SetHeight(value int) *Resize {
	t.Height = value
	return t
}

func (Resize) GetType() TransformationType {
	return TransformationTypeResize
}

func (t *Resize) Log() {
	LogTransformation(t, logrus.Fields{
		"width":  t.Width,
		"height": t.Height,
	})
}

func (t *Resize) Do(img *Image) *Image {
	img.Image = imaging.Resize(img, t.Width, t.Height, imaging.Lanczos)
	return img
}
