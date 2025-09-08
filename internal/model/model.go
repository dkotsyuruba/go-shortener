package model

type Link struct {
	Id          string `json:"id"`
	OriginalUrl string `json:"original_url"`
}

type ServerConfig struct {
	Port string
}

type ServiceConfig struct {
	BaseUrl string
}
