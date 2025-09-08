package handler_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/dkotsyuruba/go-shortener/internal/handler"
	mocks "github.com/dkotsyuruba/go-shortener/internal/handler/mocks"
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

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func readResponse(body io.Reader) string {
	responseData, _ := io.ReadAll(body)
	return string(responseData)
}
