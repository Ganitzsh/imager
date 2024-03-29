package httpv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ganitzsh/12fact/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	ContentType = "Content-Type"
)

type HTTPServerV1 struct {
	Image *ImageHandler
	service.TokenUseCase
	*service.Config
}

func NewHTTPServerV1(cfg *service.Config) *HTTPServerV1 {
	if !cfg.DevMode {
		gin.SetMode(gin.ReleaseMode)
	}
	return &HTTPServerV1{
		Config: cfg,
		Image:  NewImageHandler(),
	}
}

func (s *HTTPServerV1) SetTokenUseCase(
	value service.TokenUseCase,
) *HTTPServerV1 {
	s.TokenUseCase = value
	return s
}

func (s *HTTPServerV1) ListenAndServe() error {
	var err error

	if s.TokenUseCase == nil {
		return errors.New("No token use case")
	}

	addr := fmt.Sprintf("%s:%d", s.Host, s.HTTPPort)
	logrus.Infof("Starting HTTP Server on port %d", s.HTTPPort)
	if s.TLSEnabled {
		err = http.ListenAndServeTLS(addr, s.TLSCert, s.TLSKey, s.GetHandler())
	} else {
		err = http.ListenAndServe(addr, s.GetHandler())
	}
	if err != nil {
		logrus.Errorf("Failed to start http server: %v", err)
	}
	return err
}

func (s *HTTPServerV1) GetHandler() http.Handler {
	ret := gin.New()
	ret.Use(midCors())
	ret.Use(midLogrusLogger)
	ret.Use(midCheckErrors)
	v1 := ret.Group("/api/v1")
	{
		if !s.DevMode {
			v1.Use(midValidateToken)
		}
		s.Image.Plug(v1)
	}
	return ret
}
