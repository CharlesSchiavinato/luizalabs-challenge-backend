package cache

import "github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"

type Order interface {
	SetDetailsByOrderID(modelOrderDetails *model.OrderDetails) error
	GetDetailsByOrderID(orderID int64) (*model.OrderDetails, error)
	DelDetailsByOrderID(orderID int64) error
	ClearAll() error
}
