package mock_repository

import (
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
	"github.com/stretchr/testify/mock"
)

type MockRepositoryOrder struct {
	mock.Mock
}

func (mockRepositoryOrder *MockRepositoryOrder) LegacyBulkInsert(modelUsers *model.Users, modelOrders *model.Orders, modelOrdersProducts *model.OrdersProducts) error {
	args := mockRepositoryOrder.Called()

	return args.Error(0)
}

func (mockRepositoryOrder *MockRepositoryOrder) GetDetailsByOrderID(orderID int64) (*model.OrderDetails, error) {
	args := mockRepositoryOrder.Called()

	var modelOrderDetails *model.OrderDetails

	if args.Get(0) != nil {
		modelOrderDetails = args.Get(0).(*model.OrderDetails)
	}

	return modelOrderDetails, args.Error(1)
}

func (mockRepositoryOrder *MockRepositoryOrder) ListDetailsByRangeBuyDate(modelOrderRangeBuyDate *model.OrderRangeBuyDate) (*model.OrdersDetails, error) {
	args := mockRepositoryOrder.Called()

	var modelOrdersDetails *model.OrdersDetails

	if args.Get(0) != nil {
		modelOrdersDetails = args.Get(0).(*model.OrdersDetails)
	}

	return modelOrdersDetails, args.Error(1)
}

func (mockRepositoryOrder *MockRepositoryOrder) ListDetails() (*model.OrdersDetails, error) {
	args := mockRepositoryOrder.Called()

	var modelOrdersDetails *model.OrdersDetails

	if args.Get(0) != nil {
		modelOrdersDetails = args.Get(0).(*model.OrdersDetails)
	}

	return modelOrdersDetails, args.Error(1)
}
