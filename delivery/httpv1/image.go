package httpv1

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"

	"github.com/ganitzsh/12fact/service"
	"github.com/gin-gonic/gin"
)

type ImageHandler struct{}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{}
}

func (h *ImageHandler) Plug(e *gin.RouterGroup) error {
	image := e.Group("/images", midEnforceContentType("multipart/form-data"))
	{
		image.POST("/rotate", h.Rotate)
		image.POST("/blur", h.Blur)
		image.POST("/crop", h.Crop)
	}
	return nil
}

func extractFormFile(
	c *gin.Context,
	label string,
) (f io.Reader, ct, ext string, err error) {
	fileHeader, err := c.FormFile(label)
	if err != nil {
		return
	}
	f, err = fileHeader.Open()
	if err != nil {
		return
	}
	ct = fileHeader.Header.Get(ContentType)
	ext = filepath.Ext(fileHeader.Filename)
	return
}

func writeImage(c *gin.Context, f io.Reader, ct string) {
	buff := new(bytes.Buffer)
	if _, err := io.Copy(buff, f); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
	c.Writer.Header().Set(ContentType, ct)
	c.Writer.Write(buff.Bytes())
	c.Writer.WriteHeaderNow()
}

type RotatePayload struct {
	Angle     float64 `form:"angle"`
	ClockWise bool    `form:"clockWise"`
}

func (h *ImageHandler) Rotate(c *gin.Context) {
	p := RotatePayload{}
	if err := c.ShouldBind(&p); err != nil {
		c.Error(ErrInvalidInput)
		return
	}
	file, ct, ext, err := extractFormFile(c, "file")
	if err != nil {
		c.Error(err)
		return
	}
	f, err := service.SingleTransformImage(file, ext, &service.Rotate{
		Angle:     p.Angle,
		ClockWise: p.ClockWise,
	})
	if err != nil {
		c.Error(err)
		return
	}
	writeImage(c, f, ct)
}

type BlurPayload struct {
	Sigma float64 `form:"sigma"`
}

func (h *ImageHandler) Blur(c *gin.Context) {
	p := BlurPayload{}
	if err := c.Bind(&p); err != nil {
		c.Error(ErrInvalidInput)
		return
	}
	file, ct, ext, err := extractFormFile(c, "file")
	if err != nil {
		c.Error(err)
		return
	}
	f, err := service.SingleTransformImage(file, ext, &service.Blur{
		Sigma: p.Sigma,
	})
	if err != nil {
		c.Error(err)
		return
	}
	writeImage(c, f, ct)
}

type CropPayload struct {
	TopLeftX int `form:"topLeftX"`
	TopLeftY int `form:"topLeftY"`
	Width    int `form:"width"`
	Height   int `form:"height"`
}

func (h *ImageHandler) Crop(c *gin.Context) {
	p := CropPayload{}
	if err := c.Bind(&p); err != nil {
		c.Error(ErrInvalidInput)
		return
	}
	file, ct, ext, err := extractFormFile(c, "file")
	if err != nil {
		c.Error(err)
		return
	}
	f, err := service.SingleTransformImage(file, ext, &service.Crop{
		TopLeftX: p.TopLeftX,
		TopLeftY: p.TopLeftY,
		Width:    p.Width,
		Height:   p.Height,
	})
	if err != nil {
		c.Error(err)
		return
	}
	writeImage(c, f, ct)
}
