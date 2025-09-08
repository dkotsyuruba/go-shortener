package repository

import (
	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/dkotsyuruba/go-shortener/internal/repository/memory"
)

type Repository interface {
	Save(link *model.Link) error
	FindByID(id string) (*model.Link, bool)
}

func NewRepository() Repository {
	return memory.NewMemoryRepository()
}
