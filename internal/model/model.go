package model

type Link struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
}

type ServerConfig struct {
	Address string `env:"SERVER_ADDRESS"`
}

type ServiceConfig struct {
	BaseURL     string `env:"BASE_URL"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}
