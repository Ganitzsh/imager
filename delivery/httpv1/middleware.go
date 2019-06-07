package httpv1

import (
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/ganitzsh/12fact/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func midValidateToken(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if auth != "" {
		auth = c.Request.Header.Get("authorization")
	}
	if auth == "" {
		c.Error(ErrTokenInvalid)
		c.Abort()
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if err := service.ValidateToken(token); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.Next()
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

func midLogrusLogger(c *gin.Context) {
	path := c.Request.URL.Path
	start := time.Now()
	c.Next()
	end := time.Since(start)
	latency := int(math.Ceil(float64(end.Nanoseconds()) / 1000000.0))
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()
	clientUserAgent := c.Request.UserAgent()
	referer := c.Request.Referer()
	dataLength := c.Writer.Size()
	if dataLength < 0 {
		dataLength = 0
	}
	logrus.WithFields(logrus.Fields{
		"path":        path,
		"start":       start,
		"end":         end,
		"latency":     latency,
		"client_ip":   clientIP,
		"user_agent":  clientUserAgent,
		"referer":     referer,
		"data_length": dataLength,
		"status":      statusCode,
	}).Info("New HTTP request")
}
