package httpv1

import (
	"net/http"
	"strings"

	"github.com/ganitzsh/12fact/service"
	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error) {
	if _, ok := err.(*HandlerError); !ok {
		if se, ok := err.(*service.ServiceError); ok {
			err = NewHandlerError(
				se.Code,
				se.Error(),
				http.StatusInternalServerError,
			)
		} else {
			err = NewHandlerError(
				"unknown_error",
				err.Error(),
				http.StatusInternalServerError,
			)
		}
	}
	c.JSON(err.(*HandlerError).Status, err)
}

func midCheckErrors(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		handleError(c, c.Errors[0].Err)
		return
	}
	c.Next()
}

func midEnforceContentType(ct string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.Header.Get("Content-Type"), ct) {
			c.Error(ErrInvalidContentType)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
