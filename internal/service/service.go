package service

import (
	"errors"
	"fmt"

	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/dkotsyuruba/go-shortener/internal/repository"
	"github.com/dkotsyuruba/go-shortener/pkg/shortener"
)

type Service interface {
	Shorten(url string, userID string) (string, error)
	Get(id string) (string, error)
	GetAllByUserID(id string) ([]model.UserURLResponse, error)
	Ping() error
	ShortenBatch(batch []*model.BatchShortenRequest, userID string) ([]*model.BatchShortenResponse, error)
}

type service struct {
	repo      repository.Repository
	cfg       *model.ServiceConfig
	shortener shortener.ShortenerService
}

func NewService(
	repository repository.Repository,
	config *model.ServiceConfig,
	shortener shortener.ShortenerService,
) Service {
	return &service{
		repo:      repository,
		cfg:       config,
		shortener: shortener,
	}
}

func (s *service) Shorten(originalURL string, userID string) (string, error) {
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
		UUID:        userID,
	}

	link, err := s.repo.Save(newLink)

	return s.cfg.BaseURL + "/" + link.ID, err
}

func (s *service) ShortenBatch(batch []*model.BatchShortenRequest, userID string) ([]*model.BatchShortenResponse, error) {
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
			UUID:        userID,
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

func (s *service) Get(id string) (string, error) {
	link, found := s.repo.FindByID(id)
	if !found {
		return "", fmt.Errorf("no such short URL (%s)", id)
	}

	return link.OriginalURL, nil
}

func (s *service) GetAllByUserID(id string) ([]model.UserURLResponse, error) {
	links, err := s.repo.FindAllByUserID(id)
	if err != nil {
		return nil, err
	}

	response := make([]model.UserURLResponse, len(links))
	for i, url := range links {
		response[i] = model.UserURLResponse{
			ShortURL:    s.cfg.BaseURL + "/" + url.ID,
			OriginalURL: url.OriginalURL,
		}
	}

	return response, nil
}

func (s *service) Ping() error {
	return s.repo.Ping()
}
