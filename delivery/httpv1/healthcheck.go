package httpv1

import (
	"net/http"

	"github.com/ganitzsh/12fact/service"
	"github.com/gin-gonic/gin"
)

type HealthcheckResponse struct {
	Store bool `json:"store"`
}

func (s *HTTPServerV1) Healthcheck(c *gin.Context) {
	if s.TokenUseCase == nil || s.TokenUseCase.GetTokenStore() == nil {
		c.Error(service.ErrInvalidConfig)
		return
	}
	c.JSON(http.StatusOK, &HealthcheckResponse{
		Store: s.TokenUseCase.GetTokenStore().Up(),
	})
}
