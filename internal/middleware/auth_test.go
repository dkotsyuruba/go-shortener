package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dkotsyuruba/go-shortener/internal/middleware"
	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

type mockJWTManager struct {
	generateTokenFunc func(userID string) (string, error)
	validateTokenFunc func(token string) (string, error)
}

func (m *mockJWTManager) GenerateToken(userID string) (string, error) {
	return m.generateTokenFunc(userID)
}

func (m *mockJWTManager) ValidateToken(token string) (string, error) {
	return m.validateTokenFunc(token)
}

func TestNoCookie(t *testing.T) {
	logger := zaptest.NewLogger(t)

	m := &mockJWTManager{
		generateTokenFunc: func(userID string) (string, error) { return "token123", nil },
		validateTokenFunc: func(token string) (string, error) { return "user123", nil },
	}

	authMiddleware := middleware.NewAuthMiddleware(m, logger)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(model.UserIDContextKey)
		assert.Equal(t, "user123", userID)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler := authMiddleware.Authenticate(testHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestValidToken(t *testing.T) {
	logger := zaptest.NewLogger(t)

	m := &mockJWTManager{
		generateTokenFunc: func(userID string) (string, error) { return "token456", nil },
		validateTokenFunc: func(token string) (string, error) { return "user456", nil },
	}

	authMiddleware := middleware.NewAuthMiddleware(m, logger)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: middleware.CookieName, Value: "token456"})
	w := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(model.UserIDContextKey)
		assert.Equal(t, "user456", userID)
		w.WriteHeader(http.StatusOK)
	})

	handler := authMiddleware.Authenticate(testHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInvalidToken(t *testing.T) {
	logger := zaptest.NewLogger(t)

	m := &mockJWTManager{
		generateTokenFunc: func(userID string) (string, error) { return "token789", nil },
		validateTokenFunc: func(token string) (string, error) {
			if token == "invalidtoken" {
				return "", errors.New("invalid token")
			}
			return "user789", nil
		},
	}

	authMiddleware := middleware.NewAuthMiddleware(m, logger)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: middleware.CookieName, Value: "invalidtoken"})
	w := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(model.UserIDContextKey)
		assert.Equal(t, "user789", userID)
		w.WriteHeader(http.StatusOK)
	})

	handler := authMiddleware.Authenticate(testHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGenerateTokenError(t *testing.T) {
	logger := zaptest.NewLogger(t)

	m := &mockJWTManager{
		generateTokenFunc: func(userID string) (string, error) { return "", errors.New("gen error") },
		validateTokenFunc: nil,
	}

	authMiddleware := middleware.NewAuthMiddleware(m, logger)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called on error")
	})

	handler := authMiddleware.Authenticate(testHandler)
	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}
