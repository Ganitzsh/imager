package main

import (
	"io"
	"os"
	"path/filepath"

	trans "github.com/ganitzsh/12fact/service"
)

func main() {
	image := "./image.png"
	ext := filepath.Ext(image)
	f, err := os.Open(image)
	if err != nil {
		panic(err)
	}

	// You can now start applying transformations
	rotation := trans.NewRotate().SetAngle(90).SetClockWise(true)
	ret, err := trans.SingleTransformImage(f, ext, rotation)
	if err != nil {
		panic(err)
	}

	blur := trans.NewBlur().SetSigma(4.5)
	ret, err = trans.SingleTransformImage(ret, ext, blur)
	if err != nil {
		panic(err)
	}

	out, err := os.Create("./out.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if _, e := io.Copy(out, ret); e != nil {
		panic(e)
	}

	// You can also apply multiple transforamtions at once
	transformations := []trans.Transformation{
		trans.NewCrop().
			SetTopLeftX(40).SetTopLeftY(150).SetWidth(400).SetHeight(400),
		trans.NewRotate().SetAngle(45).SetClockWise(true),
		trans.NewBlur().SetSigma(3.7),
	}

	f.Seek(0, 0)
	ret2, err := trans.TransformImage(f, ext, transformations)
	if err != nil {
		panic(err)
	}

	out2, err := os.Create("./out2.png")
	if err != nil {
		panic(err)
	}
	defer out2.Close()

	if _, err := io.Copy(out2, ret2); err != nil {
		panic(err)
	}
}
