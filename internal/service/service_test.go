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
	mockRepo.On("Save", mock.AnythingOfType("*model.Link")).Return(&model.Link{
		ID:          id,
		OriginalURL: originalURL,
	}, nil)
	mockShortener := new(mocks.MockShortener)
	mockShortener.On("GenerateID").Return(id)

	s := service.NewService(mockRepo, config, mockShortener)

	shortenedURL, err := s.Shorten(originalURL, "")
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
	mockRepo.On("Save", mock.AnythingOfType("*model.Link")).Return(&model.Link{
		ID:          id,
		OriginalURL: originalURL,
	}, errors.New("database failure"))
	mockShortener := new(mocks.MockShortener)
	mockShortener.On("GenerateID").Return(id)

	s := service.NewService(mockRepo, config, mockShortener)

	shortenedURL, err := s.Shorten(originalURL, "")
	require.Error(t, err)
	expectedShortenedURL := config.BaseURL + "/" + id
	assert.Equal(t, expectedShortenedURL, shortenedURL)
}

func TestShortenBatchSuccess(t *testing.T) {
	config := &model.ServiceConfig{
		BaseURL: "http://example.com",
	}

	batch := []*model.BatchShortenRequest{
		{CorrelationID: "1", OriginalURL: "https://google.com"},
		{CorrelationID: "2", OriginalURL: "https://github.com"},
	}

	mockRepo := new(mocks.MockRepository)
	mockShortener := new(mocks.MockShortener)

	mockShortener.On("GenerateID").Return("id1").Once()
	mockShortener.On("GenerateID").Return("id2").Once()

	mockRepo.On("SaveAll", mock.AnythingOfType("[]*model.Link")).Return(nil)

	s := service.NewService(mockRepo, config, mockShortener)

	responses, err := s.ShortenBatch(batch, "")
	require.NoError(t, err)
	require.Len(t, responses, 2)

	assert.Equal(t, "1", responses[0].CorrelationID)
	assert.Equal(t, config.BaseURL+"/id1", responses[0].ShortURL)

	assert.Equal(t, "2", responses[1].CorrelationID)
	assert.Equal(t, config.BaseURL+"/id2", responses[1].ShortURL)
}

func TestShortenBatchFailure(t *testing.T) {
	config := &model.ServiceConfig{
		BaseURL: "http://example.com",
	}

	batch := []*model.BatchShortenRequest{
		{CorrelationID: "1", OriginalURL: "https://google.com"},
		{CorrelationID: "2", OriginalURL: "https://github.com"},
	}

	mockRepo := new(mocks.MockRepository)
	mockShortener := new(mocks.MockShortener)

	mockShortener.On("GenerateID").Return("id1").Once()
	mockShortener.On("GenerateID").Return("id2").Once()

	mockRepo.On("SaveAll", mock.AnythingOfType("[]*model.Link")).Return(errors.New("save all failure"))

	s := service.NewService(mockRepo, config, mockShortener)

	responses, err := s.ShortenBatch(batch, "")
	require.Error(t, err)
	assert.Nil(t, responses)
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

func TestGetAllByUserIDSuccess(t *testing.T) {
	config := &model.ServiceConfig{
		BaseURL: "http://example.com",
	}

	userID := "user1"
	links := []*model.Link{
		{ID: "id1", OriginalURL: "https://google.com"},
		{ID: "id2", OriginalURL: "https://github.com"},
	}

	mockRepo := new(mocks.MockRepository)
	mockShortener := new(mocks.MockShortener)

	mockRepo.On("FindAllByUserID", userID).Return(links, nil)

	s := service.NewService(mockRepo, config, mockShortener)

	responses, err := s.GetAllByUserID(userID)
	require.NoError(t, err)
	require.Len(t, responses, 2)

	assert.Equal(t, config.BaseURL+"/id1", responses[0].ShortURL)
	assert.Equal(t, "https://google.com", responses[0].OriginalURL)

	assert.Equal(t, config.BaseURL+"/id2", responses[1].ShortURL)
	assert.Equal(t, "https://github.com", responses[1].OriginalURL)
}

func TestGetAllByUserIDFailure(t *testing.T) {
	config := &model.ServiceConfig{
		BaseURL: "http://example.com",
	}

	userID := "user1"

	mockRepo := new(mocks.MockRepository)
	mockShortener := new(mocks.MockShortener)

	mockRepo.On("FindAllByUserID", userID).Return(nil, errors.New("find failure"))

	s := service.NewService(mockRepo, config, mockShortener)

	responses, err := s.GetAllByUserID(userID)
	require.Error(t, err)
	assert.Nil(t, responses)
}

func TestPingSuccess(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockShortener := new(mocks.MockShortener)

	mockRepo.On("Ping").Return(nil)

	s := service.NewService(mockRepo, &model.ServiceConfig{}, mockShortener)

	err := s.Ping()
	require.NoError(t, err)
}

func TestPingFailure(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	mockShortener := new(mocks.MockShortener)

	mockRepo.On("Ping").Return(errors.New("ping failure"))

	s := service.NewService(mockRepo, &model.ServiceConfig{}, mockShortener)

	err := s.Ping()
	require.Error(t, err)
}
