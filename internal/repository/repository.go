package repository

import (
	"github.com/dkotsyuruba/go-shortener/internal/config"
	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/dkotsyuruba/go-shortener/internal/repository/memory"
	"github.com/dkotsyuruba/go-shortener/internal/repository/postgres"
)

type Repository interface {
	Migrate() error
	Save(link *model.Link) (*model.Link, error)
	SaveAll(links []*model.Link) error
	FindByID(id string) (*model.Link, bool)
	FindByOriginalURL(id string) (*model.Link, bool)
	Ping() error
	Close() error
}

func NewRepository(config *config.Config) (Repository, error) {
	if config.DataBase.DSN != "" {
		return postgres.NewPostgresRepository(config.DataBase.DSN)
	} else {
		return memory.NewMemoryRepository(config.Service.FileStorage)
	}
}
