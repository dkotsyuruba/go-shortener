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

func (mr *MockRepository) FindById(id string) (*model.Link, bool) {
	args := mr.Called(id)
	return args.Get(0).(*model.Link), args.Bool(1)
}
