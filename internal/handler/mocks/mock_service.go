package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (fs *MockService) Shorten(url string) (string, error) {
	args := fs.Called(url)
	return args.String(0), args.Error(1)
}

func (fs *MockService) Get(id string) (string, error) {
	args := fs.Called(id)
	return args.String(0), args.Error(1)
}
