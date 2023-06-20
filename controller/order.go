package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository"
	logger "github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/logger"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/usecase"
	"github.com/hashicorp/go-hclog"
)

type Order struct {
	Title        string
	Log          hclog.Logger
	UsecaseOrder usecase.Order
}

func NewOrder(log hclog.Logger, usecaseOrder usecase.Order) *Order {
	return &Order{
		Title:        "Order",
		Log:          log,
		UsecaseOrder: usecaseOrder,
	}
}

// LegacyImport godoc
// @Summary      Importar Legado
// @Description  Importação de pedidos do sistema legado.<br/><br/>
// @Description  <strong>ATENÇÃO:</strong><br/>
// @Description  A API mantém apenas os pedidos do último arquivo importado.<br/>
// @Description  Todo o conteúdo do arquivo será desconsiderado caso ocorra algum erro durante a importação.<br/>
// @Description  Por padrão a API armazena as informações do arquivo importado em memória, portanto se a API for reiniciada os pedidos do último arquivo importado serão perdido e precisará ser importado novamente.<br/>
// @Description  É possível configurar a API para realizar o armazenamento do último arquivo importado em banco de dados.
// @Tags         Pedidos
// @Accept       json
// @Produce      json
// @Param        file   formData      file  false  "Arquivo a ser importado (formato TXT com posição fixa)" example(data_1.txt) validate(required)
// @Success      200  {object}  model.LegacyImportResult
// @Failure      400  {object}  model.Error
// @Failure      500  {object}  model.Error
// @Router       /order/legacy/import [post]
func (controllerOrder *Order) LegacyImport(rw http.ResponseWriter, req *http.Request) {
	// Parse the multipart form
	err := req.ParseMultipartForm(10 << 20) // Limit the maximum file size to 10MB

	if err != nil {
		responseError := model.BadRequestFormParsing()

		logger.LogErrorRequest(controllerOrder.Log, req, responseError.Message, err)

		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	file, fileHeader, err := req.FormFile("file")

	if err != nil {
		responseError := model.BadRequestRetrievingFile()

		logger.LogErrorRequest(controllerOrder.Log, req, responseError.Message, err)

		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	defer file.Close()

	if fileHeader.Header.Get("Content-Type") != "text/plain" {
		responseError := model.BadRequestFileType()

		logger.LogErrorRequest(controllerOrder.Log, req, responseError.Message, err)

		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	modelLegacyImportResult, err := controllerOrder.UsecaseOrder.LegacyImport(file, false)

	if err != nil {
		var responseError *model.Error

		if _, ok := err.(usecase.ErrRecordValidate); ok {
			responseError = model.BadRequestFileRecordValidate(err.Error())

			rw.WriteHeader(http.StatusBadRequest)
		} else if _, ok := err.(repository.ErrDuplicateKey); ok {
			responseError = model.BadRequestRepositoryPersist(controllerOrder.Title, err.Error())

			rw.WriteHeader(http.StatusBadRequest)
		} else {
			responseError = model.InternalServerErrorRepositoryPersist(controllerOrder.Title)

			logger.LogErrorRequest(controllerOrder.Log, req, responseError.Message, err)

			rw.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(rw).Encode(responseError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(modelLegacyImportResult)
}

// GetDetailsByOrderID godoc
// @Summary      Consultar Pedido por ID
// @Description  Retorna as informações do Pedido referente ao ID informado.
// @Tags         Pedidos
// @Accept       json
// @Produce      json
// @Param        id   path      string  false  "Número do Pedido" example(1) validate(required)
// @Success      200  {object}  model.OrderDetails
// @Failure      400  {object}  model.Error
// @Failure      404  {object}  model.Error
// @Failure      500  {object}  model.Error
// @Router       /order/{id} [get]
func (controllerOrder *Order) GetDetailsByOrderID(rw http.ResponseWriter, req *http.Request) {
	paramOrderID := strings.Split(req.URL.Path, "/")[3]

	orderID, err := strconv.ParseInt(paramOrderID, 10, 64)

	if err != nil {
		responseError := model.BadRequestParamValidate("ID invalid")
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(responseError)
		return
	}

	modelOrdersDetails, err := controllerOrder.UsecaseOrder.GetDetailsByOrderID(orderID)

	if err != nil {
		var responseError *model.Error

		if _, ok := err.(repository.ErrNotFound); ok {
			responseError = model.NotFound(controllerOrder.Title)

			rw.WriteHeader(http.StatusNotFound)
		} else {
			responseError = model.InternalServerErrorRepositoryLoad(controllerOrder.Title)

			logger.LogErrorRequest(controllerOrder.Log, req, responseError.Message, err)

			rw.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(rw).Encode(responseError)
		return
	}

	json.NewEncoder(rw).Encode(modelOrdersDetails)
}

// ListDetails godoc
// @Summary      Listar Pedidos
// @Description  Retorna todos os Pedidos ou os Pedidos referente ao período informado. O período não pode ser superior a 31 dias.
// @Tags         Pedidos
// @Accept       json
// @Produce      json
// @Param        from query      string  false  "Data da Compra Inicial (AAAA-MM-DD)" example("2020-05-23")
// @Param        to   query      string  false  "Data da Compra Final (AAAA-MM-DD)" example("2020-05-23")
// @Success      200  {object}  model.OrdersDetails
// @Failure      400  {object}  model.Error
// @Failure      404  {object}  model.Error
// @Failure      500  {object}  model.Error
// @Router       /order [get]
func (controllerOrder *Order) ListDetails(rw http.ResponseWriter, req *http.Request) {
	fromParam := req.URL.Query().Get("from")
	toParam := req.URL.Query().Get("to")

	modelOrdersDetails := &model.OrdersDetails{}
	modelOrderRangeBuyDate := &model.OrderRangeBuyDate{}
	var err error

	if fromParam == "" && toParam == "" {
		modelOrdersDetails, err = controllerOrder.UsecaseOrder.ListDetails()
	} else {
		modelOrderRangeBuyDate, err = validateQueryParamsOrderRangeBuyDate(fromParam, toParam)

		if err != nil {
			responseError := model.BadRequestParamValidate(err.Error())

			logger.LogErrorRequest(controllerOrder.Log, req, responseError.Message, err)

			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(responseError)
			return
		}

		modelOrdersDetails, err = controllerOrder.UsecaseOrder.ListDetailsByRangeBuyDate(modelOrderRangeBuyDate)
	}

	if err != nil {
		var responseError *model.Error

		if _, ok := err.(usecase.ErrParamValidate); ok {
			responseError = model.BadRequestParamValidate(err.Error())

			rw.WriteHeader(http.StatusBadRequest)
		} else if _, ok := err.(repository.ErrNotFound); ok {
			responseError = model.NotFound(controllerOrder.Title)

			rw.WriteHeader(http.StatusNotFound)
		} else {
			responseError = model.InternalServerErrorRepositoryLoad(controllerOrder.Title)

			logger.LogErrorRequest(controllerOrder.Log, req, responseError.Message, err)

			rw.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(rw).Encode(responseError)
		return
	}

	json.NewEncoder(rw).Encode(modelOrdersDetails)
}

func validateQueryParamsOrderRangeBuyDate(fromParam, toParam string) (*model.OrderRangeBuyDate, error) {
	modelOrderRangeBuyDate := &model.OrderRangeBuyDate{}

	messages := []string{}

	if fromParam == "" {
		messages = append(messages, usecase.OrderRangeBuyDateErrorMessageFromEmpty)
	} else {
		referenceDateFrom, err := time.Parse("2006-01-02", fromParam)

		if err != nil {
			messages = append(messages, usecase.OrderRangeBuyDateErrorMessageFromInvalid)
		}

		modelOrderRangeBuyDate.From = referenceDateFrom
	}

	if toParam == "" {
		messages = append(messages, usecase.OrderRangeBuyDateErrorMessageToEmpty)
	} else {
		referenceDateTo, err := time.Parse("2006-01-02", toParam)

		if err != nil {
			messages = append(messages, usecase.OrderRangeBuyDateErrorMessageToInvalid)
		}

		modelOrderRangeBuyDate.To = referenceDateTo
	}

	if len(messages) > 0 {
		return nil, errors.New(strings.Join(messages, ";"))
	}

	return modelOrderRangeBuyDate, nil
}
