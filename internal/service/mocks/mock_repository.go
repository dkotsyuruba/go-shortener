package service_test

import (
	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mr *MockRepository) Save(link *model.Link) (*model.Link, error) {
	args := mr.Called(link)
	return args.Get(0).(*model.Link), args.Error(1)
}

func (mr *MockRepository) SaveAll(link []*model.Link) error {
	args := mr.Called(link)
	return args.Error(0)
}

func (mr *MockRepository) FindByID(id string) (*model.Link, bool) {
	args := mr.Called(id)
	return args.Get(0).(*model.Link), args.Bool(1)
}

func (mr *MockRepository) FindByOriginalURL(url string) (*model.Link, bool) {
	args := mr.Called(url)
	return args.Get(0).(*model.Link), args.Bool(1)
}

func (mr *MockRepository) Close() error {
	args := mr.Called()
	return args.Error(0)
}

func (mr *MockRepository) Ping() error {
	args := mr.Called()
	return args.Error(0)
}
