package repository

import (
	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/dkotsyuruba/go-shortener/internal/repository/memory"
)

type Repository interface {
	Save(link *model.Link) error
	FindByID(id string) (*model.Link, bool)
	Persist(filename string) error
	LoadFromFile(filename string) (map[string]*model.Link, error)
}

func NewRepository(filename string) Repository {
	return memory.NewMemoryRepository(filename)
}
