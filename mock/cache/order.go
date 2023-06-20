package mock_cache

import (
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
	"github.com/stretchr/testify/mock"
)

type MockCacheOrder struct {
	mock.Mock
}

func (mockCacheOrder *MockCacheOrder) SetDetailsByOrderID(modelOrderDetails *model.OrderDetails) error {
	args := mockCacheOrder.Called()

	return args.Error(0)
}

func (mockCacheOrder *MockCacheOrder) GetDetailsByOrderID(orderID int64) (*model.OrderDetails, error) {
	args := mockCacheOrder.Called()

	var modelOrderDetails *model.OrderDetails

	if args.Get(0) != nil {
		modelOrderDetails = args.Get(0).(*model.OrderDetails)
	}

	return modelOrderDetails, args.Error(1)
}

func (mockCacheOrder *MockCacheOrder) DelDetailsByOrderID(orderID int64) error {
	args := mockCacheOrder.Called()

	return args.Error(0)
}

func (mockCacheOrder *MockCacheOrder) ClearAll() error {
	args := mockCacheOrder.Called()

	return args.Error(0)
}
