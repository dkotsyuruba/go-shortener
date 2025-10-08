package service_test

import (
	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mr *MockRepository) Save(link *model.Link) error {
	args := mr.Called(link)
	return args.Error(0)
}

func (mr *MockRepository) FindByID(id string) (*model.Link, bool) {
	args := mr.Called(id)
	return args.Get(0).(*model.Link), args.Bool(1)
}

func (mr *MockRepository) Persist(filename string) error {
	args := mr.Called(filename)
	return args.Error(0)
}

func (mr *MockRepository) LoadFromFile(filename string) (map[string]*model.Link, error) {
	args := mr.Called(filename)
	return args.Get(0).(map[string]*model.Link), args.Error(1)
}
