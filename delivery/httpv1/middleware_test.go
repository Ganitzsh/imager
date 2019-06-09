package httpv1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMidValidateToken(t *testing.T) {
	h := gin.New()
	h.Use(midCheckErrors)
	h.GET("/tmp", midValidateToken)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/tmp", nil)
	require.NoError(t, err)
	h.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Result().StatusCode)

	rec = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/tmp", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer 123")
	h.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
}
