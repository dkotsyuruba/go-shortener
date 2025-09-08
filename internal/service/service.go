package service

import (
	"fmt"

	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/dkotsyuruba/go-shortener/internal/repository"
	"github.com/dkotsyuruba/go-shortener/pkg/shortener"
)

type Service struct {
	repo      repository.Repository
	cfg       *model.ServiceConfig
	shortener shortener.ShortenerService
}

func NewService(
	repository repository.Repository,
	config *model.ServiceConfig,
	shortener shortener.ShortenerService,
) *Service {
	return &Service{
		repo:      repository,
		cfg:       config,
		shortener: shortener,
	}
}

func (s *Service) Shorten(originalURL string) (string, error) {
	id := s.shortener.GenerateID()

	newLink := &model.Link{
		ID:          id,
		OriginalURL: originalURL,
	}

	err := s.repo.Save(newLink)
	if err != nil {
		return "", err
	}

	return s.cfg.BaseURL + "/" + id, nil
}

func (s *Service) Get(id string) (string, error) {
	link, found := s.repo.FindByID(id)
	if !found {
		return "", fmt.Errorf("no such short URL (%s)", id)
	}

	return link.OriginalURL, nil
}
