package service

import (
	"errors"
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
	if originalURL == "" {
		return "", errors.New("empty URL")
	}

	id, err := s.shortener.GenerateID()
	if err != nil {
		return "", err
	}

	newLink := &model.Link{
		ID:          id,
		OriginalURL: originalURL,
	}

	link, err := s.repo.Save(newLink)

	return s.cfg.BaseURL + "/" + link.ID, err
}

func (s *Service) ShortenBatch(batch []*model.BatchShortenRequest) ([]*model.BatchShortenResponse, error) {
	modelLinks := make([]*model.Link, len(batch))
	results := make([]*model.BatchShortenResponse, len(batch))

	for idx, item := range batch {
		id, err := s.shortener.GenerateID()
		if err != nil {
			return nil, err
		}

		link := &model.Link{
			ID:          id,
			OriginalURL: item.OriginalURL,
		}

		results[idx] = &model.BatchShortenResponse{
			CorrelationID: item.CorrelationID,
			ShortURL:      s.cfg.BaseURL + "/" + link.ID,
		}

		modelLinks[idx] = link
	}

	err := s.repo.SaveAll(modelLinks)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Service) Get(id string) (string, error) {
	link, found := s.repo.FindByID(id)
	if !found {
		return "", fmt.Errorf("no such short URL (%s)", id)
	}

	return link.OriginalURL, nil
}

func (s *Service) Ping() error {
	err := s.repo.Ping()
	if err != nil {
		return err
	}

	return nil
}
