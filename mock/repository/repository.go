package mock_repository

import (
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mockRepository *MockRepository) Order() repository.Order {
	args := mockRepository.Called()
	return args.Get(0).(repository.Order)
}

func (mockRepository *MockRepository) Check() error {
	args := mockRepository.Called()

	return args.Error(0)
}

func (mockRepository *MockRepository) Close() error {
	args := mockRepository.Called()

	return args.Error(0)
}
