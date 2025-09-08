package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/dkotsyuruba/go-shortener/internal/service"
	mocks "github.com/dkotsyuruba/go-shortener/internal/service/mocks"
)

func TestShortenSuccess(t *testing.T) {
	config := &model.ServiceConfig{
		BaseURL: "http://example.com",
	}

	originalURL := "https://google.com/search?q=test"
	id := "abc123"

	mockRepo := new(mocks.MockRepository)
	mockRepo.On("Save", mock.AnythingOfType("*model.Link")).Return(nil)
	mockShortener := new(mocks.MockShortener)
	mockShortener.On("GenerateID").Return(id)

	s := service.NewService(mockRepo, config, mockShortener)

	shortenedURL, err := s.Shorten(originalURL)
	require.NoError(t, err)

	expectedShortenedURL := config.BaseURL + "/" + id
	assert.Equal(t, expectedShortenedURL, shortenedURL)
}

func TestShortenFailure(t *testing.T) {
	config := &model.ServiceConfig{
		BaseURL: "http://example.com",
	}

	originalURL := "https://google.com/search?q=test"
	id := "abc123"

	mockRepo := new(mocks.MockRepository)
	mockRepo.On("Save", mock.AnythingOfType("*model.Link")).Return(errors.New("database failure"))
	mockShortener := new(mocks.MockShortener)
	mockShortener.On("GenerateID").Return(id)

	s := service.NewService(mockRepo, config, mockShortener)

	shortenedURL, err := s.Shorten(originalURL)
	require.Error(t, err)
	assert.Empty(t, shortenedURL)
}

func TestGetSuccess(t *testing.T) {
	config := &model.ServiceConfig{
		BaseURL: "http://example.com",
	}

	id := "abc123"
	originalURL := "https://google.com/search?q=test"

	mockRepo := new(mocks.MockRepository)
	mockRepo.On("FindByID", id).Return(&model.Link{
		ID:          id,
		OriginalURL: originalURL,
	}, true)

	mockShortener := new(mocks.MockShortener)
	mockShortener.On("GenerateID").Return(id)

	s := service.NewService(mockRepo, config, mockShortener)

	actualURL, err := s.Get(id)
	require.NoError(t, err)
	assert.Equal(t, originalURL, actualURL)
}

func TestGetNotFound(t *testing.T) {
	config := &model.ServiceConfig{
		BaseURL: "http://example.com",
	}

	id := "nonexistent-id"

	mockRepo := new(mocks.MockRepository)
	mockRepo.On("FindByID", id).Return((*model.Link)(nil), false)
	mockShortener := new(mocks.MockShortener)
	mockShortener.On("GenerateID").Return(id)

	s := service.NewService(mockRepo, config, mockShortener)

	actualURL, err := s.Get(id)
	require.Error(t, err)
	assert.Empty(t, actualURL)
}
