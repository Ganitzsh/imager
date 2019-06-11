package service

import (
	"bytes"
	"errors"
	"image"
	"io"
	"sync"

	"github.com/disintegration/imaging"
)

var mut sync.Mutex
var exts = []string{
	".jpg", ".jpeg", ".png", ".tif", ".tiff", ".bmp", ".gif",
}

var formatConverter = map[string]imaging.Format{
	exts[0]: imaging.JPEG,
	exts[1]: imaging.JPEG,
	exts[2]: imaging.PNG,
	exts[3]: imaging.TIFF,
	exts[4]: imaging.TIFF,
	exts[5]: imaging.BMP,
	exts[6]: imaging.GIF,
}

var compatMatrix = map[string]bool{
	exts[0]: true,
	exts[1]: true,
	exts[2]: true,
	exts[3]: true,
	exts[4]: true,
	exts[5]: true,
	exts[6]: true,
}

type Transformation interface {
	GetType() TransformationType
	Do(img *Image) *Image
	Log()
}

func GetFormatFromExtension(ext string) imaging.Format {
	mut.Lock()
	ret := formatConverter[ext]
	mut.Unlock()
	return ret
}

func IsCompatible(ext string) bool {
	mut.Lock()
	ret := compatMatrix[ext]
	mut.Unlock()
	return ret
}

func TransformImage(
	f io.Reader, ext string,
	transformations []Transformation,
) (io.Reader, error) {
	img, err := NewImage(f, ext)
	if err != nil {
		return nil, err
	}
	for _, t := range transformations {
		if t != nil {
			t.Log()
			if err := t.Do(img).Err(); err != nil {
				return nil, err
			}
		}
	}
	return imageBuffer(img.Image, img.Format)
}

func SingleTransformImage(
	f io.Reader, ext string,
	transformation Transformation,
) (io.Reader, error) {
	if transformation == nil {
		return nil, errors.New("Invalid transformation")
	}
	return TransformImage(f, ext, []Transformation{transformation})
}

func imageBuffer(img image.Image, format imaging.Format) (io.Reader, error) {
	ret := bytes.Buffer{}
	if err := imaging.Encode(&ret, img, format); err != nil {
		return nil, err
	}
	return &ret, nil
}
