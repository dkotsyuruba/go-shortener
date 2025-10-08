package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/dkotsyuruba/go-shortener/internal/handler"
	mocks "github.com/dkotsyuruba/go-shortener/internal/handler/mocks"
	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/dkotsyuruba/go-shortener/internal/repository"
	"github.com/dkotsyuruba/go-shortener/internal/service"
	"github.com/dkotsyuruba/go-shortener/pkg/shortener"
)

func TestShortenSuccess(t *testing.T) {
	mockService := new(mocks.MockService)
	handler := handler.NewHandler(mockService)

	reqBody := []byte("https://example.com/test-url")
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))

	mockService.On("Shorten", "https://example.com/test-url").Return("short-url", nil)

	handler.Shorten(recorder, request)
	resp := recorder.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
	assert.Equal(t, "short-url", readResponse(resp.Body))
}

func TestShortenFailure(t *testing.T) {
	mockService := new(mocks.MockService)
	handler := handler.NewHandler(mockService)

	reqBody := []byte("https://example.com/test-url")
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))

	mockService.On("Shorten", "https://example.com/test-url").Return("", errors.New("oops"))

	handler.Shorten(recorder, request)
	resp := recorder.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetSuccess(t *testing.T) {
	mockService := new(mocks.MockService)
	handler := handler.NewHandler(mockService)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/abcdef", nil)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "abcdef")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

	mockService.On("Get", "abcdef").Return("https://example.com/test-url", nil)

	handler.Get(recorder, request)
	resp := recorder.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "https://example.com/test-url", resp.Header.Get("Location"))
}

func TestGetFailure(t *testing.T) {
	mockService := new(mocks.MockService)
	handler := handler.NewHandler(mockService)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/nonexistent-id", nil)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "nonexistent-id")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

	mockService.On("Get", "nonexistent-id").Return("", errors.New("not found"))

	handler.Get(recorder, request)
	resp := recorder.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestShortenAPI(t *testing.T) {
	mockRepo := repository.NewRepository("test.json")
	cfg := &model.ServiceConfig{
		BaseURL: "http://localhost:8080",
	}

	shortenerService := shortener.NewRealShortenerService()
	srv := service.NewService(mockRepo, cfg, shortenerService)
	handlers := handler.NewHandler(srv)

	body := `{"url": "https://practicum.yandex.ru"}`
	reader := bytes.NewBufferString(body)
	req := httptest.NewRequest("POST", "/api/shorten", reader)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handlers.ShortenJSON(recorder, req)

	res := recorder.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var resp = model.Response{}
	err := json.NewDecoder(res.Body).Decode(&resp)

	defer res.Body.Close()
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Result)
}

func readResponse(body io.Reader) string {
	responseData, _ := io.ReadAll(body)
	return string(responseData)
}
