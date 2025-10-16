package mocks

import (
	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Shorten(url string) (string, error) {
	args := m.Called(url)
	return "short-url", args.Error(1)
}

func (m *MockService) ShortenBatch(batch []*model.BatchShortenRequest) ([]*model.BatchShortenResponse, error) {
	args := m.Called(batch)
	return args.Get(0).([]*model.BatchShortenResponse), args.Error(1)
}

func (m *MockService) Get(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func (m *MockService) Ping() error {
	args := m.Called()
	return args.Error(0)
}
