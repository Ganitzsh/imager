package httpv1_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/ganitzsh/12fact/delivery/httpv1"
	"github.com/ganitzsh/12fact/service"
	"github.com/stretchr/testify/require"
)

func createPNGFormFile(w *multipart.Writer, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", "image/png")
	return w.CreatePart(h)
}

func getDevServer() *httpv1.HTTPServerV1 {
	return httpv1.NewHTTPServerV1(&service.Config{
		MaxImageSize: 25165824,
		DevMode:      true,
	})
}

func TestRotateValid(t *testing.T) {
	srv := getDevServer()
	path := "../../example/img.png"
	file, err := os.Open(path)
	defer file.Close()
	require.NoError(t, err)
	rec := httptest.NewRecorder()
	buff := &bytes.Buffer{}
	body := multipart.NewWriter(buff)
	part, err := createPNGFormFile(body, filepath.Base(path))
	require.NoError(t, err)
	_, err = io.Copy(part, file)
	require.NoError(t, err)
	require.NoError(t, body.WriteField("angle", "45"))
	require.NoError(t, body.WriteField("clockWise", "true"))
	require.NoError(t, body.Close())

	req, err := http.NewRequest(http.MethodPost, "/api/v1/images/rotate", buff)
	require.NoError(t, err)
	req.Header.Set("Content-Type", body.FormDataContentType())
	srv.GetHandler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Equal(t, "image/png", rec.Result().Header.Get("Content-Type"))
}

func TestBlurValid(t *testing.T) {
	srv := getDevServer()
	path := "../../example/img.png"
	file, err := os.Open(path)
	defer file.Close()
	require.NoError(t, err)
	rec := httptest.NewRecorder()
	buff := &bytes.Buffer{}
	body := multipart.NewWriter(buff)
	part, err := createPNGFormFile(body, filepath.Base(path))
	require.NoError(t, err)
	_, err = io.Copy(part, file)
	require.NoError(t, err)
	require.NoError(t, body.WriteField("sigma", "5.76"))
	require.NoError(t, body.Close())

	req, err := http.NewRequest(http.MethodPost, "/api/v1/images/blur", buff)
	require.NoError(t, err)
	req.Header.Set("Content-Type", body.FormDataContentType())
	srv.GetHandler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Equal(t, "image/png", rec.Result().Header.Get("Content-Type"))
}

func TestCropValid(t *testing.T) {
	srv := getDevServer()
	path := "../../example/img.png"
	file, err := os.Open(path)
	defer file.Close()
	require.NoError(t, err)
	rec := httptest.NewRecorder()
	buff := &bytes.Buffer{}
	body := multipart.NewWriter(buff)
	part, err := createPNGFormFile(body, filepath.Base(path))
	require.NoError(t, err)
	_, err = io.Copy(part, file)
	require.NoError(t, err)
	require.NoError(t, body.WriteField("topLeftX", "30"))
	require.NoError(t, body.WriteField("topLeftY", "30"))
	require.NoError(t, body.WriteField("width", "300"))
	require.NoError(t, body.WriteField("height", "300"))
	require.NoError(t, body.Close())

	req, err := http.NewRequest(http.MethodPost, "/api/v1/images/crop", buff)
	require.NoError(t, err)
	req.Header.Set("Content-Type", body.FormDataContentType())
	srv.GetHandler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Equal(t, "image/png", rec.Result().Header.Get("Content-Type"))
}
