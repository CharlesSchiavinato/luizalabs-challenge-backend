package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"reflect"
	"testing"

	mock_usecase "github.com/CharlesSchiavinato/luizalabs-challenge-backend/mock/usecase"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/usecase"
	"github.com/hashicorp/go-hclog"
)

type MockFile struct {
	io.Reader
}

func (f *MockFile) Close() error {
	// Implement the Close method if needed
	return nil
}

func TestOrderLegacyImport(t *testing.T) {
	modelLegacyImportResult := model.LegacyImportResult{
		Users:    10,
		Orders:   11,
		Products: 20,
	}

	type test struct {
		name        string
		reqFormData func() (*multipart.Writer, *bytes.Buffer)
		resBody     interface{}
		wantResCode int
		wantResBody interface{}
		mockOn      func(*mock_usecase.MockUsecaseOrder)
	}

	tests := []test{
		{
			name: "FormParsingError",
			reqFormData: func() (*multipart.Writer, *bytes.Buffer) {
				// Create a new test file with content
				fileContent :=
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308"
				fileBuffer := bytes.NewBufferString(fileContent)

				// Create a new HTTP request with a file upload
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				fileWriter, _ := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file2"; filename="file.txt"`},
					"Content-Type":        []string{"text/csv"},
				})

				if _, err := io.Copy(fileWriter, fileBuffer); err != nil {
					t.Fatalf("Failed to write file content to form file: %v", err)
				}
				writer.Close()

				return writer, &bytes.Buffer{}
			},
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestFormParsing(),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
			},
		},
		{
			name: "FormFileRetrievingError",
			reqFormData: func() (*multipart.Writer, *bytes.Buffer) {
				// Create a new test file with content
				fileContent :=
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308"
				fileBuffer := bytes.NewBufferString(fileContent)

				// Create a new HTTP request with a file upload
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				fileWriter, _ := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file2"; filename="file.txt"`},
					"Content-Type":        []string{"text/csv"},
				})

				if _, err := io.Copy(fileWriter, fileBuffer); err != nil {
					t.Fatalf("Failed to write file content to form file: %v", err)
				}
				writer.Close()

				return writer, body
			},
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestRetrievingFile(),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
			},
		},
		{
			name: "FormFileContentTypeError",
			reqFormData: func() (*multipart.Writer, *bytes.Buffer) {
				// Create a new test file with content
				fileContent :=
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308"
				fileBuffer := bytes.NewBufferString(fileContent)

				// Create a new HTTP request with a file upload
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				fileWriter, _ := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="file.txt"`},
					"Content-Type":        []string{"text/csv"},
				})

				if _, err := io.Copy(fileWriter, fileBuffer); err != nil {
					t.Fatalf("Failed to write file content to form file: %v", err)
				}
				writer.Close()

				return writer, body
			},
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestFileType(),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
			},
		},
		{
			name: "BadRequestFileRecordValidate",
			reqFormData: func() (*multipart.Writer, *bytes.Buffer) {
				// Create a new test file with content
				fileContent :=
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308"
				fileBuffer := bytes.NewBufferString(fileContent)

				// Create a new HTTP request with a file upload
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				fileWriter, _ := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="file.txt"`},
					"Content-Type":        []string{"text/plain"},
				})

				if _, err := io.Copy(fileWriter, fileBuffer); err != nil {
					t.Fatalf("Failed to write file content to form file: %v", err)
				}
				writer.Close()

				return writer, body
			},
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestFileRecordValidate("Record Error"),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("LegacyImport").Return(nil, usecase.ErrRecordValidate{Message: "Record Error"})
			},
		},
		{
			name: "BadRequestDuplicateKeyError",
			reqFormData: func() (*multipart.Writer, *bytes.Buffer) {
				// Create a new test file with content
				fileContent :=
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308"
				fileBuffer := bytes.NewBufferString(fileContent)

				// Create a new HTTP request with a file upload
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				fileWriter, _ := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="file.txt"`},
					"Content-Type":        []string{"text/plain"},
				})

				if _, err := io.Copy(fileWriter, fileBuffer); err != nil {
					t.Fatalf("Failed to write file content to form file: %v", err)
				}
				writer.Close()

				return writer, body
			},
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestRepositoryPersist("Order", "Duplicate Key"),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("LegacyImport").Return(nil, repository.ErrDuplicateKey{Message: "Duplicate Key"})
			},
		},
		{
			name: "InternalServerError",
			reqFormData: func() (*multipart.Writer, *bytes.Buffer) {
				// Create a new test file with content
				fileContent :=
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308"
				fileBuffer := bytes.NewBufferString(fileContent)

				// Create a new HTTP request with a file upload
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				fileWriter, _ := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="file.txt"`},
					"Content-Type":        []string{"text/plain"},
				})

				if _, err := io.Copy(fileWriter, fileBuffer); err != nil {
					t.Fatalf("Failed to write file content to form file: %v", err)
				}
				writer.Close()

				return writer, body
			},
			resBody:     &model.Error{},
			wantResCode: http.StatusInternalServerError,
			wantResBody: model.InternalServerErrorRepositoryPersist("Order"),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("LegacyImport").Return(nil, errors.New("InternalServerError"))
			},
		},
		{
			name: "Success",
			reqFormData: func() (*multipart.Writer, *bytes.Buffer) {
				// Create a new test file with content
				fileContent :=
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308"
				fileBuffer := bytes.NewBufferString(fileContent)

				// Create a new HTTP request with a file upload
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				fileWriter, _ := writer.CreatePart(textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="file.txt"`},
					"Content-Type":        []string{"text/plain"},
				})

				if _, err := io.Copy(fileWriter, fileBuffer); err != nil {
					t.Fatalf("Failed to write file content to form file: %v", err)
				}
				writer.Close()

				return writer, body
			},
			resBody:     model.LegacyImportResult{},
			wantResCode: http.StatusOK,
			wantResBody: model.LegacyImportResult{},
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("LegacyImport").Return(&modelLegacyImportResult, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
			mockUsecaseOrder := new(mock_usecase.MockUsecaseOrder)

			tt.mockOn(mockUsecaseOrder)

			controllerOrder := NewOrder(log, mockUsecaseOrder)

			writer, body := tt.reqFormData()

			url := fmt.Sprintf("/api/order/legacy/import")

			req := httptest.NewRequest(http.MethodPost, url, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			handler := http.HandlerFunc(controllerOrder.LegacyImport)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if !reflect.DeepEqual(res.Code, tt.wantResCode) {
				t.Errorf("LegacyImport() got res.code = %v, want %v", res.Code, tt.wantResCode)
			}

			json.NewDecoder(res.Body).Decode(tt.resBody)

			if !reflect.DeepEqual(tt.resBody, tt.wantResBody) {
				t.Errorf("LegacyImport() got res.body = %v, want %v", tt.resBody, tt.wantResBody)
			}
		})
	}
}

func TestOrderGetDetailsByOrderID(t *testing.T) {
	modelOrderDetails := model.OrderDetails{
		UserID:   70,
		UserName: "Palmer Prosacco",
		Orders: []model.OrderDetailsOrder{
			{
				OrderID: 753,
				BuyDate: "2021-03-08",
				Total:   1836.74,
				Products: []model.OrderDetailsProduct{
					{
						ID:    3,
						Value: 1836.74,
					},
				},
			},
		},
	}

	type test struct {
		name        string
		reqParam    string
		resBody     interface{}
		wantResCode int
		wantResBody interface{}
		mockOn      func(*mock_usecase.MockUsecaseOrder)
	}

	tests := []test{
		{
			name:        "RequestParamError",
			reqParam:    "X",
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate("ID invalid"),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
			},
		},
		{
			name:        "NotFoundError",
			reqParam:    "783",
			resBody:     &model.Error{},
			wantResCode: http.StatusNotFound,
			wantResBody: model.NotFound("Order"),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("GetDetailsByOrderID").Return(nil, repository.ErrNotFound{Message: "not found"})
			},
		},
		{
			name:        "InternalServerError",
			reqParam:    "783",
			resBody:     &model.Error{},
			wantResCode: http.StatusInternalServerError,
			wantResBody: model.InternalServerErrorRepositoryLoad("Order"),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("GetDetailsByOrderID").Return(nil, errors.New("InternalServerError"))
			},
		},
		{
			name:        "Success",
			reqParam:    "783",
			resBody:     modelOrderDetails,
			wantResCode: http.StatusOK,
			wantResBody: modelOrderDetails,
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("GetDetailsByOrderID").Return(&modelOrderDetails, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
			mockUsecaseOrder := new(mock_usecase.MockUsecaseOrder)

			tt.mockOn(mockUsecaseOrder)

			controllerOrder := NewOrder(log, mockUsecaseOrder)

			url := fmt.Sprintf("/api/order/%v", tt.reqParam)

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			handler := http.HandlerFunc(controllerOrder.GetDetailsByOrderID)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if !reflect.DeepEqual(res.Code, tt.wantResCode) {
				t.Errorf("LegacyImport() got res.code = %v, want %v", res.Code, tt.wantResCode)
			}

			json.NewDecoder(res.Body).Decode(tt.resBody)

			if !reflect.DeepEqual(tt.resBody, tt.wantResBody) {
				t.Errorf("LegacyImport() got res.body = %v, want %v", tt.resBody, tt.wantResBody)
			}
		})
	}
}

func TestOrderListDetails(t *testing.T) {
	modelOrdersDetails := model.OrdersDetails{
		{
			UserID:   70,
			UserName: "Palmer Prosacco",
			Orders: []model.OrderDetailsOrder{
				{
					OrderID: 753,
					BuyDate: "2021-03-08",
					Total:   1836.74,
					Products: []model.OrderDetailsProduct{
						{
							ID:    3,
							Value: 1836.74,
						},
					},
				},
			},
		},
	}

	type test struct {
		name        string
		reqParam    string
		resBody     interface{}
		wantResCode int
		wantResBody interface{}
		mockOn      func(*mock_usecase.MockUsecaseOrder)
	}

	tests := []test{
		{
			name:        "ParamFromEmptyError",
			reqParam:    "?from=&to=2020-01-01",
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate(usecase.OrderRangeBuyDateErrorMessageFromEmpty),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
			},
		},
		{
			name:        "ParamFromInvalidError",
			reqParam:    "?from=2020-13-01&to=2020-01-01",
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate(usecase.OrderRangeBuyDateErrorMessageFromInvalid),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
			},
		},
		{
			name:        "ParamToEmptyError",
			reqParam:    "?from=2020-01-01&to=",
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate(usecase.OrderRangeBuyDateErrorMessageToEmpty),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
			},
		},
		{
			name:        "ParamToInvalidError",
			reqParam:    "?from=2020-01-01&to=2020-13-01",
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate(usecase.OrderRangeBuyDateErrorMessageToInvalid),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
			},
		},
		{
			name:        "ParamValidateError",
			reqParam:    fmt.Sprintf("?from=%v&to=2020-01-01", usecase.OrderBuyDateMin.AddDate(0, 0, -1).Format("2006-01-02")),
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate(usecase.OrderRangeBuyDateErrorMessageFromBetween),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("ListDetailsByRangeBuyDate").Return(nil, usecase.ErrParamValidate{Message: usecase.OrderRangeBuyDateErrorMessageFromBetween})
			},
		},
		{
			name:        "NotFoundError",
			reqParam:    "",
			resBody:     &model.Error{},
			wantResCode: http.StatusNotFound,
			wantResBody: model.NotFound("Order"),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("ListDetails").Return(nil, repository.ErrNotFound{Message: "not found"})
			},
		},
		{
			name:        "InternalServerError",
			reqParam:    "",
			resBody:     &model.Error{},
			wantResCode: http.StatusInternalServerError,
			wantResBody: model.InternalServerErrorRepositoryLoad("Order"),
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("ListDetails").Return(nil, errors.New("InternalServerError"))
			},
		},
		{
			name:        "FromToSuccess",
			reqParam:    "?from=2020-01-01&to=2020-01-31",
			resBody:     &model.OrdersDetails{},
			wantResCode: http.StatusOK,
			wantResBody: &modelOrdersDetails,
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("ListDetailsByRangeBuyDate").Return(&modelOrdersDetails, nil)
			},
		},
		{
			name:        "AllSuccess",
			reqParam:    "",
			resBody:     &model.OrdersDetails{},
			wantResCode: http.StatusOK,
			wantResBody: &modelOrdersDetails,
			mockOn: func(mockUsecaseOrder *mock_usecase.MockUsecaseOrder) {
				mockUsecaseOrder.On("ListDetails").Return(&modelOrdersDetails, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
			mockUsecaseOrder := new(mock_usecase.MockUsecaseOrder)

			tt.mockOn(mockUsecaseOrder)

			controllerOrder := NewOrder(log, mockUsecaseOrder)

			url := fmt.Sprintf("/api/order%v", tt.reqParam)

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			handler := http.HandlerFunc(controllerOrder.ListDetails)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if !reflect.DeepEqual(res.Code, tt.wantResCode) {
				t.Errorf("LegacyImport() got res.code = %v, want %v", res.Code, tt.wantResCode)
			}

			json.NewDecoder(res.Body).Decode(tt.resBody)

			if !reflect.DeepEqual(tt.resBody, tt.wantResBody) {
				t.Errorf("LegacyImport() got res.body = %v, want %v", tt.resBody, tt.wantResBody)
			}
		})
	}
}
