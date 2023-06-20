package usecase

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/cache"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/util"
)

var (
	OrderBuyDateMin                            = time.Date(1900, 01, 01, 00, 00, 00, 000, time.UTC)
	OrderBuyDateMax                            = time.Now().UTC()
	OrderErrorMessageRecordSize                = "Record size not equal 95"
	OrderErrorMessageUserIDInvalid             = "UserID invalid"
	OrderErrorMessageUserNameInvalid           = "UserName invalid"
	OrderErrorMessageOrderIDInvalid            = "OrderID invalid"
	OrderErrorMessageProductIDInvalid          = "ProductID invalid"
	OrderErrorMessageProductValueInvalid       = "ProductValue invalid"
	OrderErrorMessageBuyDateInvalid            = "BuyDate invalid"
	OrderErrorMessageBuyDateBetween            = fmt.Sprintf("BuyDate value is not between %v and %v", OrderBuyDateMin.Format("2006-01-02"), OrderBuyDateMax.Format("2006-01-02"))
	OrderRangeBuyDateErrorMessageFromEmpty     = "The param from is empty"
	OrderRangeBuyDateErrorMessageFromInvalid   = "The param from is invalid"
	OrderRangeBuyDateErrorMessageFromBetween   = fmt.Sprintf("The param from value is not between %v and %v", OrderBuyDateMin.Format("2006-01-02"), OrderBuyDateMax.Format("2006-01-02"))
	OrderRangeBuyDateErrorMessageToEmpty       = "The param to is empty"
	OrderRangeBuyDateErrorMessageToInvalid     = "The param to is invalid"
	OrderRangeBuyDateErrorMessageToBetween     = fmt.Sprintf("The param to value is not between %v and %v", OrderBuyDateMin.Format("2006-01-02"), OrderBuyDateMax.Format("2006-01-02"))
	OrderRangeBuyDateErrorMessageToSmallerFrom = "The param to is smaller the param from"
	OrderRangeBuyDateErrorMessageRangeError    = "the range is greater than 31 days"
)

type Order interface {
	LegacyImport(file io.Reader, hasHeader bool) (*model.LegacyImportResult, error)
	GetDetailsByOrderID(orderID int64) (*model.OrderDetails, error)
	ListDetailsByRangeBuyDate(modelOrderRangeBuyDate *model.OrderRangeBuyDate) (*model.OrdersDetails, error)
	ListDetails() (*model.OrdersDetails, error)
}

type UseCaseOrder struct {
	Repository repository.Repository
	Cache      cache.Cache
}

func NewOrder(repository repository.Repository, cache cache.Cache) Order {
	return &UseCaseOrder{
		Repository: repository,
		Cache:      cache,
	}
}

func (usecaseOrder *UseCaseOrder) GetDetailsByOrderID(orderID int64) (*model.OrderDetails, error) {
	modelOrdersDetails, err := usecaseOrder.Cache.Order().GetDetailsByOrderID(orderID)

	if err == nil {
		return modelOrdersDetails, err
	}

	modelOrdersDetails, err = usecaseOrder.Repository.Order().GetDetailsByOrderID(orderID)

	if err == nil {
		usecaseOrder.Cache.Order().SetDetailsByOrderID(modelOrdersDetails)
	}

	return modelOrdersDetails, err
}

func (usecaseOrder *UseCaseOrder) ListDetails() (*model.OrdersDetails, error) {
	return usecaseOrder.Repository.Order().ListDetails()
}

func (usecaseOrder *UseCaseOrder) ListDetailsByRangeBuyDate(modelOrderRangeBuyDate *model.OrderRangeBuyDate) (*model.OrdersDetails, error) {
	err := OrderRangeBuyDateValidate(modelOrderRangeBuyDate)

	if err != nil {
		return nil, err
	}

	return usecaseOrder.Repository.Order().ListDetailsByRangeBuyDate(modelOrderRangeBuyDate)
}

func (usecaseOrder *UseCaseOrder) LegacyImport(file io.Reader, hasHeader bool) (*model.LegacyImportResult, error) {
	scanner := bufio.NewScanner(file)

	if hasHeader {
		scanner.Scan()
	}

	linesCount := 0
	modelUsers := model.Users{}
	modelOrders := model.Orders{}
	modelOrdersProducts := model.OrdersProducts{}

	modelLegacyRecordsError := model.LegacyRecordsError{}
	mapUsers := make(map[int64]int)
	mapOrders := make(map[int64]int)

	for scanner.Scan() {
		linesCount++
		record := scanner.Text()

		modelLegacy, err := recordToLegacy(record)

		if err != nil {
			modelLegacyRecordsError = append(modelLegacyRecordsError, model.LegacyRecordError{Line: int64(linesCount), Message: err.Error()})
			continue
		}

		legacyUser(modelLegacy, &modelUsers, mapUsers)
		legacyOrder(modelLegacy, &modelOrders, mapOrders)
		legacyProduct(modelLegacy, &modelOrdersProducts)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(modelLegacyRecordsError) > 0 {
		jsonBytes, _ := json.Marshal(modelLegacyRecordsError)
		return nil, ErrRecordValidate{Message: string(jsonBytes)}
	}

	usecaseOrder.Cache.Order().ClearAll()

	err := usecaseOrder.Repository.Order().LegacyBulkInsert(&modelUsers, &modelOrders, &modelOrdersProducts)

	if err != nil {
		return nil, err
	}

	modelLegacyImportResult := &model.LegacyImportResult{
		Users:    len(modelUsers),
		Orders:   len(modelOrders),
		Products: len(modelOrdersProducts),
	}
	return modelLegacyImportResult, err
}

func recordToLegacy(record string) (*model.Legacy, error) {
	if len(record) != 95 {
		return nil, errors.New(OrderErrorMessageRecordSize)
	}

	modelRecordLegacy := &model.LegacyRecord{
		UserID:       record[0:10],
		UserName:     strings.TrimSpace(record[10:55]),
		OrderID:      record[55:65],
		ProductID:    record[65:75],
		ProductValue: strings.TrimSpace(record[75:87]),
		BuyDate:      record[87:95],
	}

	modelLegacy := &model.Legacy{}
	errMessages := []string{}
	var err error

	modelLegacy.UserID, err = strconv.ParseInt(modelRecordLegacy.UserID, 10, 64)

	if err != nil {
		errMessages = append(errMessages, OrderErrorMessageUserIDInvalid)
	}

	modelLegacy.UserName = util.FormatTitle(modelRecordLegacy.UserName)

	if len(modelLegacy.UserName) < 2 {
		errMessages = append(errMessages, OrderErrorMessageUserNameInvalid)
	}

	modelLegacy.OrderID, err = strconv.ParseInt(modelRecordLegacy.OrderID, 10, 64)

	if err != nil {
		errMessages = append(errMessages, OrderErrorMessageOrderIDInvalid)
	}

	modelLegacy.ProductID, err = strconv.ParseInt(modelRecordLegacy.ProductID, 10, 64)

	if err != nil {
		errMessages = append(errMessages, OrderErrorMessageProductIDInvalid)
	}

	modelLegacy.ProductValue, err = strconv.ParseFloat(modelRecordLegacy.ProductValue, 64)

	if err != nil {
		errMessages = append(errMessages, OrderErrorMessageProductValueInvalid)
	}

	buyDate, err := time.Parse("20060102", modelRecordLegacy.BuyDate)

	if err != nil {
		errMessages = append(errMessages, OrderErrorMessageBuyDateInvalid)
	} else {
		if buyDate.Before(OrderBuyDateMin) || buyDate.After(OrderBuyDateMax) {
			errMessages = append(errMessages, OrderErrorMessageBuyDateBetween)
		}
	}

	if len(errMessages) > 0 {
		return nil, errors.New(strings.Join(errMessages, ";"))
	}

	modelLegacy.BuyDate = buyDate.Format("2006-01-02")
	modelLegacy.ImportedAt = time.Now().UTC()

	return modelLegacy, nil
}

func legacyUser(modelLegacy *model.Legacy, modelUsers *model.Users, mapUsers map[int64]int) {
	index, ok := mapUsers[modelLegacy.UserID]

	if !ok {
		*modelUsers = append(*modelUsers, model.User{
			ID:   modelLegacy.UserID,
			Name: modelLegacy.UserName,
		})

		index = len(*modelUsers) - 1
		mapUsers[modelLegacy.UserID] = index
	}

	return
}

func legacyOrder(modelLegacy *model.Legacy, modelOrders *model.Orders, mapOrders map[int64]int) {
	index, ok := mapOrders[modelLegacy.OrderID]

	if !ok {
		*modelOrders = append(*modelOrders, model.Order{
			ID:      modelLegacy.OrderID,
			UserID:  modelLegacy.UserID,
			BuyDate: modelLegacy.BuyDate,
			Total:   modelLegacy.ProductValue,
		})

		index = len(*modelOrders) - 1
		mapOrders[modelLegacy.OrderID] = index
	} else {
		if (*modelOrders)[index].UserID != modelLegacy.UserID {
			log.Fatal("UserID divergent")
		}

		(*modelOrders)[index].Total =
			util.MathRoundPrecision((*modelOrders)[index].Total+modelLegacy.ProductValue, 2)
	}

	return
}

func legacyProduct(modelLegacy *model.Legacy, modelOrdersProducts *model.OrdersProducts) {
	*modelOrdersProducts = append(*modelOrdersProducts, model.OrderProduct{
		OrderID:      modelLegacy.OrderID,
		ProductID:    modelLegacy.ProductID,
		ProductValue: modelLegacy.ProductValue,
	})
}

func OrderRangeBuyDateValidate(modelOrderRangeBuyDate *model.OrderRangeBuyDate) error {
	messages := []string{}

	if modelOrderRangeBuyDate.From.IsZero() {
		messages = append(messages, OrderRangeBuyDateErrorMessageFromEmpty)
	} else if modelOrderRangeBuyDate.From.Before(OrderBuyDateMin) ||
		modelOrderRangeBuyDate.From.After(OrderBuyDateMax) {
		messages = append(messages, OrderRangeBuyDateErrorMessageFromBetween)
	}

	if modelOrderRangeBuyDate.To.IsZero() {
		messages = append(messages, OrderRangeBuyDateErrorMessageToEmpty)
	} else if modelOrderRangeBuyDate.To.Before(OrderBuyDateMin) ||
		modelOrderRangeBuyDate.To.After(OrderBuyDateMax) {
		messages = append(messages, OrderRangeBuyDateErrorMessageToBetween)
	}

	if len(messages) > 0 {
		return ErrParamValidate{Message: strings.Join(messages, ";")}
	}

	dateDiffDays := int64(modelOrderRangeBuyDate.To.Sub(modelOrderRangeBuyDate.From).Hours() / 24)

	if dateDiffDays < 0 {
		messages = append(messages, OrderRangeBuyDateErrorMessageToSmallerFrom)
	} else if dateDiffDays > 31 {
		messages = append(messages, OrderRangeBuyDateErrorMessageRangeError)
	}

	if len(messages) > 0 {
		return ErrParamValidate{Message: strings.Join(messages, ";")}
	}

	return nil
}
