package repository

import (
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository"
)

var (
	orderModelUsers          = model.Users{}
	orderModelOrders         = model.Orders{}
	orderModelOrdersProducts = model.OrdersProducts{}
	orderMapUsers            = make(map[int64]int)
	orderMapOrders           = make(map[int64]int)
	orderMapOrdersProducts   = make(map[int64][]int)
)

type InMemoryOrder struct{}

func NewOrder() repository.Order {
	return &InMemoryOrder{}
}

func (*InMemoryOrder) LegacyBulkInsert(modelUsers *model.Users, modelOrders *model.Orders, modelOrdersProducts *model.OrdersProducts) error {
	orderModelUsers = *modelUsers
	orderModelOrders = *modelOrders
	orderModelOrdersProducts = *modelOrdersProducts

	mapUsers := make(map[int64]int)
	mapOrders := make(map[int64]int)
	mapOrdersProducts := make(map[int64][]int)

	pos := 0

	for pos < len(orderModelUsers) {
		mapUsers[orderModelUsers[pos].ID] = pos
		mapOrders[orderModelOrders[pos].ID] = pos
		mapOrdersProducts[orderModelOrdersProducts[pos].OrderID] = append(mapOrdersProducts[orderModelOrdersProducts[pos].OrderID], pos)
		pos++
	}

	for pos < len(orderModelOrders) {
		mapOrders[orderModelOrders[pos].ID] = pos
		mapOrdersProducts[orderModelOrdersProducts[pos].OrderID] = append(mapOrdersProducts[orderModelOrdersProducts[pos].OrderID], pos)
		pos++
	}

	for pos < len(orderModelOrdersProducts) {
		mapOrdersProducts[orderModelOrdersProducts[pos].OrderID] = append(mapOrdersProducts[orderModelOrdersProducts[pos].OrderID], pos)
		pos++
	}

	orderMapUsers = mapUsers
	orderMapOrders = mapOrders
	orderMapOrdersProducts = mapOrdersProducts

	return nil
}

func (inMemoryOrder *InMemoryOrder) GetDetailsByOrderID(orderID int64) (*model.OrderDetails, error) {
	orderIndex, ok := orderMapOrders[orderID]

	if !ok {
		return nil, repository.ErrNotFound{Message: "not found"}
	}

	modelOrdersDetails := model.OrdersDetails{}
	mapOrdersDetails := make(map[int64]int)

	inMemoryOrder.convertToDetails(&modelOrdersDetails, mapOrdersDetails, &orderModelOrders[orderIndex])

	return &(modelOrdersDetails)[0], nil
}

func (inMemoryOrder *InMemoryOrder) ListDetailsByRangeBuyDate(modelOrderRangeBuyDate *model.OrderRangeBuyDate) (*model.OrdersDetails, error) {
	orderRangeBuyDateFrom := modelOrderRangeBuyDate.From.Format("2006-01-02")
	orderRangeBuyDateTo := modelOrderRangeBuyDate.To.Format("2006-01-02")

	modelOrdersDetails := model.OrdersDetails{}
	mapOrdersDetails := make(map[int64]int)

	for _, modelOrder := range orderModelOrders {
		if modelOrder.BuyDate >= orderRangeBuyDateFrom && modelOrder.BuyDate <= orderRangeBuyDateTo {
			inMemoryOrder.convertToDetails(&modelOrdersDetails, mapOrdersDetails, &modelOrder)
		}
	}

	if len(modelOrdersDetails) == 0 {
		return nil, repository.ErrNotFound{Message: "not found"}
	}

	return &modelOrdersDetails, nil
}

func (inMemoryOrder *InMemoryOrder) ListDetails() (*model.OrdersDetails, error) {
	modelOrdersDetails := model.OrdersDetails{}
	mapOrdersDetails := make(map[int64]int)

	for _, modelOrder := range orderModelOrders {
		inMemoryOrder.convertToDetails(&modelOrdersDetails, mapOrdersDetails, &modelOrder)
	}

	if len(modelOrdersDetails) == 0 {
		return nil, repository.ErrNotFound{Message: "not found"}
	}

	return &modelOrdersDetails, nil
}

func (*InMemoryOrder) convertToDetails(modelOrdersDetails *model.OrdersDetails, mapOrdersDetails map[int64]int, modelOrder *model.Order) {
	userIndex := orderMapUsers[modelOrder.UserID]

	modelUser := &orderModelUsers[userIndex]

	modelOrderDetailsProducts := []model.OrderDetailsProduct{}

	for _, orderProductIndex := range orderMapOrdersProducts[modelOrder.ID] {
		modelOrderProduct := orderModelOrdersProducts[orderProductIndex]
		modelOrderDetailsProducts = append(modelOrderDetailsProducts, model.OrderDetailsProduct{
			ID:    modelOrderProduct.ProductID,
			Value: modelOrderProduct.ProductValue,
		})
	}

	modelOrderDetailsOrder := model.OrderDetailsOrder{
		OrderID:  modelOrder.ID,
		BuyDate:  modelOrder.BuyDate,
		Total:    modelOrder.Total,
		Products: modelOrderDetailsProducts,
	}

	detailsIndex, ok := mapOrdersDetails[modelOrder.UserID]

	if !ok {
		*modelOrdersDetails = append(*modelOrdersDetails, model.OrderDetails{
			UserID:   modelOrder.UserID,
			UserName: modelUser.Name,
			Orders:   []model.OrderDetailsOrder{modelOrderDetailsOrder},
		})

		mapOrdersDetails[modelOrder.UserID] = len(*modelOrdersDetails) - 1
	} else {
		(*modelOrdersDetails)[detailsIndex].Orders = append((*modelOrdersDetails)[detailsIndex].Orders, modelOrderDetailsOrder)
	}

	return
}
