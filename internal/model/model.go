package model

type Link struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
}

type ServerConfig struct {
	Port string
}

type ServiceConfig struct {
	BaseURL string
}
