package middleware

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/dkotsyuruba/go-shortener/internal/model"
	jwtpkg "github.com/dkotsyuruba/go-shortener/pkg/jwt"
)

const CookieName = "User"

type AuthMiddleware struct {
	jwtManager *jwtpkg.JWTManager
	logger     *zap.Logger
}

func NewAuthMiddleware(jwtManager *jwtpkg.JWTManager, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		logger:     logger,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string
		cookie, err := r.Cookie(CookieName)
		if err != nil || cookie.Value == "" {
			userID, err = m.generateAndSetCookie(w)
			if err != nil {
				m.logger.Error("failed to generate user id", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			userID, err = m.jwtManager.ValidateToken(cookie.Value)
			if err != nil {
				if err == jwtpkg.ErrMissingUserID {
					http.Error(w, "Unauthorized: missing user ID", http.StatusUnauthorized)
					return
				}
				userID, err = m.generateAndSetCookie(w)
				if err != nil {
					m.logger.Error("failed to generate user id after invalid token", zap.Error(err))
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
			}
		}

		ctx := context.WithValue(r.Context(), model.UserIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) generateAndSetCookie(w http.ResponseWriter) (string, error) {
	token, err := m.jwtManager.GenerateToken("")
	if err != nil {
		return "", err
	}

	userID, err := m.jwtManager.ValidateToken(token)
	if err != nil {
		return "", err
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	return userID, nil
}
