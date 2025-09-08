package config

import (
	"flag"

	"github.com/dkotsyuruba/go-shortener/internal/model"
)

type Config struct {
	Server  *model.ServerConfig
	Service *model.ServiceConfig
}

func InitConfig() *Config {
	cfg := Config{
		Server:  &model.ServerConfig{},
		Service: &model.ServiceConfig{},
	}
	cfg.LoadConfig()

	return &cfg
}

func (c *Config) LoadConfig() {
	flag.StringVar(&c.Server.Port, "a", ":8080", "HTTP server startup address")
	flag.StringVar(&c.Service.BaseUrl, "b", "http://localhost:8080", "Base URL for shortened URL")
	flag.Parse()
}
