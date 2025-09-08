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

func (s *Service) Shorten(originalUrl string) (string, error) {
	id := s.shortener.GenerateId()

	newLink := &model.Link{
		Id:          id,
		OriginalUrl: originalUrl,
	}

	err := s.repo.Save(newLink)
	if err != nil {
		return "", err
	}

	return s.cfg.BaseUrl + "/" + id, nil
}

func (s *Service) Get(id string) (string, error) {
	link, found := s.repo.FindById(id)
	if !found {
		return "", fmt.Errorf("no such short Url (%s)", id)
	}

	return link.OriginalUrl, nil
}
