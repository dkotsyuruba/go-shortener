package config

import (
	"flag"

	"github.com/caarlos0/env"
	"github.com/dkotsyuruba/go-shortener/internal/model"
)

type Config struct {
	Server   *model.ServerConfig
	Service  *model.ServiceConfig
	DataBase *model.DatabaseConfig
}

func InitConfig() *Config {
	cfg := Config{
		Server:   &model.ServerConfig{},
		Service:  &model.ServiceConfig{},
		DataBase: &model.DatabaseConfig{},
	}

	cfg.LoadCommandLineConfig()
	cfg.LoadEnvConfig()

	return &cfg
}

func (c *Config) LoadEnvConfig() error {
	if err := env.Parse(c.Server); err != nil {
		return err
	}

	if err := env.Parse(c.Service); err != nil {
		return err
	}

	if err := env.Parse(c.DataBase); err != nil {
		return err
	}

	return nil
}

func (c *Config) LoadCommandLineConfig() {
	flag.StringVar(&c.Server.Address, "a", ":8080", "HTTP server startup address")
	flag.StringVar(&c.Service.BaseURL, "b", "http://localhost:8080", "Base URL for shortened URLs")
	flag.StringVar(&c.Service.FileStorage, "f", "storage.json", "File storage for shortened URLs")
	flag.StringVar(&c.DataBase.DSN, "d", "", "PostgreSQL connection string")
	flag.Parse()
}
