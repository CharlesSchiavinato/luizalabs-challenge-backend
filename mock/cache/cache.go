package mock_cache

import (
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/cache"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (mockCache *MockCache) Order() cache.Order {
	args := mockCache.Called()
	return args.Get(0).(cache.Order)
}

func (mockCache *MockCache) Check() error {
	args := mockCache.Called()

	return args.Error(0)
}

func (mockCache *MockCache) Close() error {
	args := mockCache.Called()

	return args.Error(0)
}
