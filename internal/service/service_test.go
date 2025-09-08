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
		BaseUrl: "http://example.com",
	}

	originalUrl := "https://google.com/search?q=test"
	id := "abc123"

	mockRepo := new(mocks.MockRepository)
	mockRepo.On("Save", mock.AnythingOfType("*model.Link")).Return(nil)
	mockShortener := new(mocks.MockShortener)
	mockShortener.On("GenerateId").Return(id)

	s := service.NewService(mockRepo, config, mockShortener)

	shortenedUrl, err := s.Shorten(originalUrl)
	require.NoError(t, err)

	expectedShortenedUrl := config.BaseUrl + "/" + id
	assert.Equal(t, expectedShortenedUrl, shortenedUrl)
}

func TestShortenFailure(t *testing.T) {
	config := &model.ServiceConfig{
		BaseUrl: "http://example.com",
	}

	originalUrl := "https://google.com/search?q=test"
	id := "abc123"

	mockRepo := new(mocks.MockRepository)
	mockRepo.On("Save", mock.AnythingOfType("*model.Link")).Return(errors.New("database failure"))
	mockShortener := new(mocks.MockShortener)
	mockShortener.On("GenerateId").Return(id)

	s := service.NewService(mockRepo, config, mockShortener)

	shortenedUrl, err := s.Shorten(originalUrl)
	require.Error(t, err)
	assert.Empty(t, shortenedUrl)
}

func TestGetSuccess(t *testing.T) {
	config := &model.ServiceConfig{
		BaseUrl: "http://example.com",
	}

	id := "abc123"
	originalUrl := "https://google.com/search?q=test"

	mockRepo := new(mocks.MockRepository)
	mockRepo.On("FindById", id).Return(&model.Link{
		Id:          id,
		OriginalUrl: originalUrl,
	}, true)

	mockShortener := new(mocks.MockShortener)
	mockShortener.On("GenerateId").Return(id)

	s := service.NewService(mockRepo, config, mockShortener)

	actualUrl, err := s.Get(id)
	require.NoError(t, err)
	assert.Equal(t, originalUrl, actualUrl)
}

func TestGetNotFound(t *testing.T) {
	config := &model.ServiceConfig{
		BaseUrl: "http://example.com",
	}

	id := "nonexistent-id"

	mockRepo := new(mocks.MockRepository)
	mockRepo.On("FindById", id).Return((*model.Link)(nil), false)
	mockShortener := new(mocks.MockShortener)
	mockShortener.On("GenerateId").Return(id)

	s := service.NewService(mockRepo, config, mockShortener)

	actualUrl, err := s.Get(id)
	require.Error(t, err)
	assert.Empty(t, actualUrl)
}
