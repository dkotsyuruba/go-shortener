package model

import "errors"

type Link struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
	UUID        string `json:"user_id"`
}

type ServerConfig struct {
	Address string `env:"SERVER_ADDRESS"`
}

type ServiceConfig struct {
	BaseURL     string `env:"BASE_URL"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
}

type DatabaseConfig struct {
	DSN string `env:"DATABASE_DSN"`
}

type AuthConfig struct {
	SecretKey string `env:"SECRET_KEY"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type BatchShortenRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchShortenResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UserURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type contextKey string

var UserIDContextKey = contextKey("user_id")

var ErrDuplicatedURL = errors.New("URL already exists")
var ErrNotFound = errors.New("URL nor found")
