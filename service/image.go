package service

import (
	"image"
	"io"

	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
)

type Image struct {
	image.Image
	err    error
	Format imaging.Format
}

func (i *Image) Err() error {
	return i.err
}

func NewImage(f io.Reader, ext string) (*Image, error) {
	if !IsCompatible(ext) {
		return nil, ErrUnsupportedFormat
	}
	img, err := imaging.Decode(f)
	if err != nil {
		logrus.Errorf("failed to decode image: %v", err)
		return nil, err
	}
	return &Image{
		Image:  img,
		Format: GetFormatFromExtension(ext),
	}, nil
}

func (i *Image) SetFormat(value imaging.Format) *Image {
	i.Format = value
	return i
}
