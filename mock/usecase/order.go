package mock_usecase

import (
	"io"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
	"github.com/stretchr/testify/mock"
)

type MockUsecaseOrder struct {
	mock.Mock
}

func (mockUsecaseOrder *MockUsecaseOrder) GetDetailsByOrderID(orderID int64) (*model.OrderDetails, error) {
	args := mockUsecaseOrder.Called()

	var modelOrderDetails *model.OrderDetails

	if args.Get(0) != nil {
		modelOrderDetails = args.Get(0).(*model.OrderDetails)
	}

	return modelOrderDetails, args.Error(1)
}

func (mockUsecaseOrder *MockUsecaseOrder) LegacyImport(file io.Reader, hasHeader bool) (*model.LegacyImportResult, error) {
	args := mockUsecaseOrder.Called()

	var modelLegacyImportResult *model.LegacyImportResult

	if args.Get(0) != nil {
		modelLegacyImportResult = args.Get(0).(*model.LegacyImportResult)
	}

	return modelLegacyImportResult, args.Error(1)
}

func (mockUsecaseOrder *MockUsecaseOrder) ListDetailsByRangeBuyDate(modelOrderRangeBuyDate *model.OrderRangeBuyDate) (*model.OrdersDetails, error) {
	args := mockUsecaseOrder.Called()

	var modelOrdersDetails *model.OrdersDetails

	if args.Get(0) != nil {
		modelOrdersDetails = args.Get(0).(*model.OrdersDetails)
	}

	return modelOrdersDetails, args.Error(1)
}

func (mockUsecaseOrder *MockUsecaseOrder) ListDetails() (*model.OrdersDetails, error) {
	args := mockUsecaseOrder.Called()

	var modelOrdersDetails *model.OrdersDetails

	if args.Get(0) != nil {
		modelOrdersDetails = args.Get(0).(*model.OrdersDetails)
	}

	return modelOrdersDetails, args.Error(1)
}
