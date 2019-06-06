package httpv1

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"

	pb "github.com/ganitzsh/12fact/proto"
	"github.com/ganitzsh/12fact/service"
)

type ImageHandler struct {
	Engine *gin.Engine
}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{}
}

func (h *ImageHandler) GetEngine() http.Handler {
	return nil
}

func (h *ImageHandler) Rotate(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		c.Error(err)
		return
	}
	any, err := ptypes.MarshalAny(&pb.RotateImageRequest{})
	if err != nil {
		c.Error(err)
		return
	}
	service.TransformImageFunc(f, filepath.Ext(f.Filename))
}
