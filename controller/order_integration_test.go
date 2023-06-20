package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"reflect"
	"strings"
	"testing"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
	cache "github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/cache/redis"
	repository "github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository/in_memory"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/usecase"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/util"
	"github.com/hashicorp/go-hclog"
)

var (
	testIntegrationLog                  = hclog.New(&hclog.LoggerOptions{Level: hclog.LevelFromString("OFF")})
	testIntegrationConfig, _            = util.LoadConfig("./../")
	testIntegrationRepository, _        = repository.NewInMemory(testIntegrationConfig)
	testIntegrationCache, _             = cache.NewRedis(testIntegrationConfig)
	testIntegrationUsecaseOrder         = usecase.NewOrder(testIntegrationRepository, testIntegrationCache)
	testIntegrationControllerOrder      = NewOrder(testIntegrationLog, testIntegrationUsecaseOrder)
	testIntegrationControllerOrderTitle = "Order"
)

func TestIntegrationOrder(t *testing.T) {
	testIntegrationOrderLegacyImport(t)
	testIntegrationOrderGetDetailsByOrderID(t)
	testIntegrationOrderListDetails(t)
}

func testIntegrationOrderLegacyImport(t *testing.T) {
	type test struct {
		name        string
		reqParam    string
		reqFormData func() (*multipart.Writer, *bytes.Buffer)
		resBody     interface{}
		wantResCode int
		wantResBody func() interface{}
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
			wantResBody: func() interface{} {
				return model.BadRequestFormParsing()
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
			wantResBody: func() interface{} {
				return model.BadRequestRetrievingFile()
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
			wantResBody: func() interface{} {
				return model.BadRequestFileType()
			},
		},
		{
			name: "BadRequestFileRecordValidate",
			reqFormData: func() (*multipart.Writer, *bytes.Buffer) {
				// Create a new test file with content
				fileContent :=
					"000000007x                                            o000000075x000000000x     1836.7x20211308"
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
			wantResBody: func() interface{} {
				jsonBytes, _ := json.Marshal(model.LegacyRecordsError{
					{
						Line: 1,
						Message: strings.Join([]string{
							usecase.OrderErrorMessageUserIDInvalid,
							usecase.OrderErrorMessageUserNameInvalid,
							usecase.OrderErrorMessageOrderIDInvalid,
							usecase.OrderErrorMessageProductIDInvalid,
							usecase.OrderErrorMessageProductValueInvalid,
							usecase.OrderErrorMessageBuyDateInvalid,
						},
							";"),
					},
				})

				return model.BadRequestFileRecordValidate(string(jsonBytes))
			},
		},
		{
			name: "Success",
			reqFormData: func() (*multipart.Writer, *bytes.Buffer) {
				// Create a new test file with content
				fileContent := []string{
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308",
					"0000000075                                  Bobbie Batz00000007980000000002     1578.5720211116",
					"0000000075                                  Bobbie Batz00000005230000000003      586.7420210903",
					"0000000070                              Palmer Prosacco00000007530000000003     1009.5420210308",
				}

				fileBuffer := bytes.NewBufferString(strings.Join(fileContent, "\n"))

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
			resBody:     &model.LegacyImportResult{},
			wantResCode: http.StatusOK,
			wantResBody: func() interface{} {
				return &model.LegacyImportResult{Users: 2, Orders: 3, Products: 4}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer, body := tt.reqFormData()

			url := fmt.Sprintf("/api/order/legacy/import")

			req := httptest.NewRequest(http.MethodPost, url, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			handler := http.HandlerFunc(testIntegrationControllerOrder.LegacyImport)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if !reflect.DeepEqual(res.Code, tt.wantResCode) {
				t.Errorf("LegacyImport() got res.code = %v, want %v", res.Code, tt.wantResCode)
			}

			json.NewDecoder(res.Body).Decode(tt.resBody)

			wantResBody := tt.wantResBody()

			if !reflect.DeepEqual(tt.resBody, wantResBody) {
				t.Errorf("LegacyImport() got res.body = %v, want %v", tt.resBody, wantResBody)
			}
		})
	}
}

func testIntegrationOrderGetDetailsByOrderID(t *testing.T) {
	type test struct {
		name        string
		reqParam    string
		resBody     interface{}
		wantResCode int
		wantResBody interface{}
	}

	tests := []test{
		{
			name:        "RequestParamError",
			reqParam:    "X",
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate("ID invalid"),
		},
		{
			name:        "NotFoundError",
			reqParam:    "70",
			resBody:     &model.Error{},
			wantResCode: http.StatusNotFound,
			wantResBody: model.NotFound("Order"),
		},
		{
			name:        "Success",
			reqParam:    "753",
			resBody:     &model.OrderDetails{},
			wantResCode: http.StatusOK,
			wantResBody: &model.OrderDetails{
				UserID:   70,
				UserName: "Palmer Prosacco",
				Orders: []model.OrderDetailsOrder{
					{
						OrderID: 753,
						BuyDate: "2021-03-08",
						Total:   2846.28,
						Products: []model.OrderDetailsProduct{
							{
								ID:    3,
								Value: 1836.74,
							},
							{
								ID:    3,
								Value: 1009.54,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/order/%v", tt.reqParam)

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			handler := http.HandlerFunc(testIntegrationControllerOrder.GetDetailsByOrderID)
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

func testIntegrationOrderListDetails(t *testing.T) {
	type test struct {
		name        string
		reqParam    string
		resBody     interface{}
		wantResCode int
		wantResBody interface{}
	}

	tests := []test{
		{
			name:        "ParamFromEmptyError",
			reqParam:    "?from=&to=2020-01-01",
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate(usecase.OrderRangeBuyDateErrorMessageFromEmpty),
		},
		{
			name:        "ParamFromInvalidError",
			reqParam:    "?from=2020-13-01&to=2020-01-01",
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate(usecase.OrderRangeBuyDateErrorMessageFromInvalid),
		},
		{
			name:        "ParamToEmptyError",
			reqParam:    "?from=2020-01-01&to=",
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate(usecase.OrderRangeBuyDateErrorMessageToEmpty),
		},
		{
			name:        "ParamToInvalidError",
			reqParam:    "?from=2020-01-01&to=2020-13-01",
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate(usecase.OrderRangeBuyDateErrorMessageToInvalid),
		},
		{
			name:        "ParamValidateError",
			reqParam:    fmt.Sprintf("?from=%v&to=2020-01-01", usecase.OrderBuyDateMin.AddDate(0, 0, -1).Format("2006-01-02")),
			resBody:     &model.Error{},
			wantResCode: http.StatusBadRequest,
			wantResBody: model.BadRequestParamValidate(usecase.OrderRangeBuyDateErrorMessageFromBetween),
		},
		{
			name:        "NotFoundError",
			reqParam:    "?from=2020-01-01&to=2020-01-01",
			resBody:     &model.Error{},
			wantResCode: http.StatusNotFound,
			wantResBody: model.NotFound("Order"),
		},
		{
			name:        "FromToSuccess",
			reqParam:    "?from=2021-11-16&to=2021-11-16",
			resBody:     &model.OrdersDetails{},
			wantResCode: http.StatusOK,
			wantResBody: &model.OrdersDetails{
				{
					UserID:   75,
					UserName: "Bobbie Batz",
					Orders: []model.OrderDetailsOrder{
						{
							OrderID: 798,
							BuyDate: "2021-11-16",
							Total:   1578.57,
							Products: []model.OrderDetailsProduct{
								{
									ID:    2,
									Value: 1578.57,
								},
							},
						},
					},
				},
			},
		},
		{
			name:        "AllSuccess",
			reqParam:    "",
			resBody:     &model.OrdersDetails{},
			wantResCode: http.StatusOK,
			wantResBody: &model.OrdersDetails{
				{
					UserID:   70,
					UserName: "Palmer Prosacco",
					Orders: []model.OrderDetailsOrder{
						{
							OrderID: 753,
							BuyDate: "2021-03-08",
							Total:   2846.28,
							Products: []model.OrderDetailsProduct{
								{
									ID:    3,
									Value: 1836.74,
								},
								{
									ID:    3,
									Value: 1009.54,
								},
							},
						},
					},
				},
				{
					UserID:   75,
					UserName: "Bobbie Batz",
					Orders: []model.OrderDetailsOrder{
						{
							OrderID: 798,
							BuyDate: "2021-11-16",
							Total:   1578.57,
							Products: []model.OrderDetailsProduct{
								{
									ID:    2,
									Value: 1578.57,
								},
							},
						},
						{
							OrderID: 523,
							BuyDate: "2021-09-03",
							Total:   586.74,
							Products: []model.OrderDetailsProduct{
								{
									ID:    3,
									Value: 586.74,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/order%v", tt.reqParam)

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			handler := http.HandlerFunc(testIntegrationControllerOrder.ListDetails)
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
