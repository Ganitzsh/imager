package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
}
