package repository

import "github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"

type Order interface {
	LegacyBulkInsert(modelUsers *model.Users, modelOrders *model.Orders, modelOrdersProducts *model.OrdersProducts) error
	GetDetailsByOrderID(orderID int64) (*model.OrderDetails, error)
	ListDetailsByRangeBuyDate(modelOrderRangeBuyDate *model.OrderRangeBuyDate) (*model.OrdersDetails, error)
	ListDetails() (*model.OrdersDetails, error)
}
