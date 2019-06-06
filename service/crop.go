package service

import (
	"image"

	"github.com/disintegration/imaging"
)

type Crop struct {
	Transformation

	TopLeftX int
	TopLeftY int
	Width    int
	Height   int
}

func NewCrop() *Crop {
	return &Crop{}
}

func (t *Crop) SetTopLeftX(value int) *Crop {
	t.TopLeftX = value
	return t
}

func (t *Crop) SetTopLeftY(value int) *Crop {
	t.TopLeftY = value
	return t
}

func (t *Crop) SetWidth(value int) *Crop {
	t.Width = value
	return t
}

func (t *Crop) SetHeight(value int) *Crop {
	t.Height = value
	return t
}

func (Crop) GetType() TransformationType {
	return TransformationTypeCrop
}

func (t *Crop) Do(img *Image) *Image {
	img.Image = imaging.Crop(img, image.Rect(
		t.TopLeftX, t.TopLeftY,
		t.TopLeftX+t.Width, t.TopLeftY+t.Height,
	))
	return img
}
