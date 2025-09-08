package service_test

import "github.com/stretchr/testify/mock"

type MockShortener struct {
	mock.Mock
}

func (m *MockShortener) GenerateID() (string, error) {
	args := m.Called()
	return args.String(0), nil
}
