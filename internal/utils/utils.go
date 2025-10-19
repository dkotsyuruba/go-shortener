package utils

import (
	"net/http"

	"github.com/dkotsyuruba/go-shortener/internal/model"
)

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(model.UserIDContextKey).(string)
	return userID, ok
}
