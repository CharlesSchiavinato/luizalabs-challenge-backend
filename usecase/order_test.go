package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	mock_cache "github.com/CharlesSchiavinato/luizalabs-challenge-backend/mock/cache"
	mock_repository "github.com/CharlesSchiavinato/luizalabs-challenge-backend/mock/repository"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
)

func TestOrderLegacyImport(t *testing.T) {
	type test struct {
		name           string
		inputFile      func() io.Reader
		inputHasHeader bool
		wantResult     *model.LegacyImportResult
		wantError      func() error
		mockOn         func(*mock_repository.MockRepository, *mock_cache.MockCache)
	}

	tests := []test{
		{
			name: "RecordSizeError",
			inputFile: func() io.Reader {
				lines := []string{
					"0000000070                             Palmer Prosacco00000007530000000003     1836.7420210308",
				}

				file := bytes.NewBufferString(strings.Join(lines, "\n"))
				return file
			},
			inputHasHeader: false,
			wantResult:     nil,
			wantError: func() error {
				jsonBytes, _ := json.Marshal(model.LegacyRecordsError{
					{
						Line:    1,
						Message: OrderErrorMessageRecordSize,
					},
				})

				err := ErrRecordValidate{Message: string(jsonBytes)}
				return err
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name: "UserIDInvalidError",
			inputFile: func() io.Reader {
				lines := []string{
					"000000007x                              Palmer Prosacco00000007530000000003     1836.7420210308",
				}

				file := bytes.NewBufferString(strings.Join(lines, "\n"))
				return file
			},
			inputHasHeader: false,
			wantResult:     nil,
			wantError: func() error {
				jsonBytes, _ := json.Marshal(model.LegacyRecordsError{
					{
						Line:    1,
						Message: OrderErrorMessageUserIDInvalid,
					},
				})

				err := ErrRecordValidate{Message: string(jsonBytes)}
				return err
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name: "UserNameInvalidError",
			inputFile: func() io.Reader {
				lines := []string{
					"0000000070                                            o00000007530000000003     1836.7420210308",
				}

				file := bytes.NewBufferString(strings.Join(lines, "\n"))
				return file
			},
			inputHasHeader: false,
			wantResult:     nil,
			wantError: func() error {
				jsonBytes, _ := json.Marshal(model.LegacyRecordsError{
					{
						Line:    1,
						Message: OrderErrorMessageUserNameInvalid,
					},
				})

				err := ErrRecordValidate{Message: string(jsonBytes)}
				return err
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name: "OrderIDInvalidError",
			inputFile: func() io.Reader {
				lines := []string{
					"0000000070                              Palmer Prosacco000000075x0000000003     1836.7420210308",
				}

				file := bytes.NewBufferString(strings.Join(lines, "\n"))
				return file
			},
			inputHasHeader: false,
			wantResult:     nil,
			wantError: func() error {
				jsonBytes, _ := json.Marshal(model.LegacyRecordsError{
					{
						Line:    1,
						Message: OrderErrorMessageOrderIDInvalid,
					},
				})

				err := ErrRecordValidate{Message: string(jsonBytes)}
				return err
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name: "ProductIDInvalidError",
			inputFile: func() io.Reader {
				lines := []string{
					"0000000070                              Palmer Prosacco0000000753000000000x     1836.7420210308",
				}

				file := bytes.NewBufferString(strings.Join(lines, "\n"))
				return file
			},
			inputHasHeader: false,
			wantResult:     nil,
			wantError: func() error {
				jsonBytes, _ := json.Marshal(model.LegacyRecordsError{
					{
						Line:    1,
						Message: OrderErrorMessageProductIDInvalid,
					},
				})

				err := ErrRecordValidate{Message: string(jsonBytes)}
				return err
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name: "ProductValueInvalidError",
			inputFile: func() io.Reader {
				lines := []string{
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7x20210308",
				}

				file := bytes.NewBufferString(strings.Join(lines, "\n"))
				return file
			},
			inputHasHeader: false,
			wantResult:     nil,
			wantError: func() error {
				jsonBytes, _ := json.Marshal(model.LegacyRecordsError{
					{
						Line:    1,
						Message: OrderErrorMessageProductValueInvalid,
					},
				})

				err := ErrRecordValidate{Message: string(jsonBytes)}
				return err
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name: "BuyDateInvalidError",
			inputFile: func() io.Reader {
				lines := []string{
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420211308",
				}

				file := bytes.NewBufferString(strings.Join(lines, "\n"))
				return file
			},
			inputHasHeader: false,
			wantResult:     nil,
			wantError: func() error {
				jsonBytes, _ := json.Marshal(model.LegacyRecordsError{
					{
						Line:    1,
						Message: OrderErrorMessageBuyDateInvalid,
					},
				})

				err := ErrRecordValidate{Message: string(jsonBytes)}
				return err
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name: "BuyDateBetweenError",
			inputFile: func() io.Reader {
				lines := []string{
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7418990308",
				}

				file := bytes.NewBufferString(strings.Join(lines, "\n"))
				return file
			},
			inputHasHeader: false,
			wantResult:     nil,
			wantError: func() error {
				jsonBytes, _ := json.Marshal(model.LegacyRecordsError{
					{
						Line:    1,
						Message: OrderErrorMessageBuyDateBetween,
					},
				})

				err := ErrRecordValidate{Message: string(jsonBytes)}
				return err
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name: "RepositoryError",
			inputFile: func() io.Reader {
				lines := []string{
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308",
					"0000000075                                  Bobbie Batz00000007980000000002     1578.5720211116",
					"0000000049                               Ken Wintheiser00000005230000000003      586.7420210903",
					"0000000014                                 Clelia Hills00000001460000000001      673.4920211125",
					"0000000057                          Elidia Gulgowski IV00000006200000000000     1417.2520210919",
					"0000000080                                 Tabitha Kuhn00000008770000000003      817.1320210612",
					"0000000023                                  Logan Lynch00000002530000000002      322.1220210523",
					"0000000015                                   Bonny Koss00000001530000000004        80.820210701",
					"0000000017                              Ethan Langworth00000001690000000000      865.1820210409",
					"0000000077                         Mrs. Stephen Trantow00000008440000000005     1288.7720211127",
					"0000000061                           Dimple Bergstrom I00000006710000000004       43.3620211104",
					"0000000077                         Mrs. Stephen Trantow00000008320000000006      961.3720210513",
					"0000000041                           Dr. Dexter Rolfson00000004470000000003     1563.4720210630",
					"0000000078                                    Wade Mraz00000008610000000003      224.9720210910",
					"0000000002                           Augustus Aufderhar00000000220000000000       190.820210530",
					"0000000025                             Frederica Cremin00000002760000000004      113.7520211103",
					"0000000069                             Dr. Tyree Rogahn00000007430000000002      1401.620210317",
					"0000000001                              Sammie Baumbach00000000070000000002       96.4720210528",
					"0000000077                         Mrs. Stephen Trantow00000008480000000004      1689.020210325",
					"0000000075                                  Bobbie Batz00000008120000000001      707.9620211103",
				}

				return bytes.NewBufferString(strings.Join(lines, "\n"))
			},
			inputHasHeader: true,
			wantResult:     nil,
			wantError: func() error {
				return errors.New("LegacyBulkInsert Error")
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				mockRepositoryOrder := new(mock_repository.MockRepositoryOrder)
				mockRepositoryOrder.On("LegacyBulkInsert").Return(errors.New("LegacyBulkInsert Error"))
				mockRepository.On("Order").Return(mockRepositoryOrder)

				mockCacheOrder := new(mock_cache.MockCacheOrder)
				mockCacheOrder.On("ClearAll").Return(nil)
				mockCache.On("Order").Return(mockCacheOrder)
			},
		},
		{
			name: "HasHeaderError",
			inputFile: func() io.Reader {
				lines := []string{
					"    USERID                                     USERNAME   ORDERID PRODUCTIDPRODUCTVALUE BUYDATE",
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308",
					"0000000075                                  Bobbie Batz00000007980000000002     1578.5720211116",
					"0000000049                               Ken Wintheiser00000005230000000003      586.7420210903",
					"0000000014                                 Clelia Hills00000001460000000001      673.4920211125",
					"0000000057                          Elidia Gulgowski IV00000006200000000000     1417.2520210919",
					"0000000080                                 Tabitha Kuhn00000008770000000003      817.1320210612",
					"0000000023                                  Logan Lynch00000002530000000002      322.1220210523",
					"0000000015                                   Bonny Koss00000001530000000004        80.820210701",
					"0000000017                              Ethan Langworth00000001690000000000      865.1820210409",
					"0000000077                         Mrs. Stephen Trantow00000008440000000005     1288.7720211127",
					"0000000061                           Dimple Bergstrom I00000006710000000004       43.3620211104",
					"0000000077                         Mrs. Stephen Trantow00000008320000000006      961.3720210513",
					"0000000041                           Dr. Dexter Rolfson00000004470000000003     1563.4720210630",
					"0000000078                                    Wade Mraz00000008610000000003      224.9720210910",
					"0000000002                           Augustus Aufderhar00000000220000000000       190.820210530",
					"0000000025                             Frederica Cremin00000002760000000004      113.7520211103",
					"0000000069                             Dr. Tyree Rogahn00000007430000000002      1401.620210317",
					"0000000001                              Sammie Baumbach00000000070000000002       96.4720210528",
					"0000000077                         Mrs. Stephen Trantow00000008480000000004      1689.020210325",
					"0000000075                                  Bobbie Batz00000008120000000001      707.9620211103",
					"0000000070                              Palmer Prosacco00000007530000000003     1009.5420210308",
				}

				return bytes.NewBufferString(strings.Join(lines, "\n"))
			},
			inputHasHeader: false,
			wantResult:     nil,
			wantError: func() error {
				jsonBytes, _ := json.Marshal(model.LegacyRecordsError{
					{
						Line: 1,
						Message: strings.Join(
							[]string{
								OrderErrorMessageUserIDInvalid,
								OrderErrorMessageOrderIDInvalid,
								OrderErrorMessageProductIDInvalid,
								OrderErrorMessageProductValueInvalid,
								OrderErrorMessageBuyDateInvalid,
							},
							";",
						),
					},
				})

				err := ErrRecordValidate{Message: string(jsonBytes)}
				return err
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				mockRepositoryOrder := new(mock_repository.MockRepositoryOrder)
				mockRepositoryOrder.On("LegacyBulkInsert").Return(nil)
				mockRepository.On("Order").Return(mockRepositoryOrder)

				mockCacheOrder := new(mock_cache.MockCacheOrder)
				mockCacheOrder.On("ClearAll").Return(nil)
				mockCache.On("Order").Return(mockCacheOrder)
			},
		},
		{
			name: "HasHeaderSuccess",
			inputFile: func() io.Reader {
				lines := []string{
					"    USERID                                     USERNAME   ORDERID PRODUCTIDPRODUCTVALUE BUYDATE",
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308",
					"0000000075                                  Bobbie Batz00000007980000000002     1578.5720211116",
					"0000000049                               Ken Wintheiser00000005230000000003      586.7420210903",
					"0000000014                                 Clelia Hills00000001460000000001      673.4920211125",
					"0000000057                          Elidia Gulgowski IV00000006200000000000     1417.2520210919",
					"0000000080                                 Tabitha Kuhn00000008770000000003      817.1320210612",
					"0000000023                                  Logan Lynch00000002530000000002      322.1220210523",
					"0000000015                                   Bonny Koss00000001530000000004        80.820210701",
					"0000000017                              Ethan Langworth00000001690000000000      865.1820210409",
					"0000000077                         Mrs. Stephen Trantow00000008440000000005     1288.7720211127",
					"0000000061                           Dimple Bergstrom I00000006710000000004       43.3620211104",
					"0000000077                         Mrs. Stephen Trantow00000008320000000006      961.3720210513",
					"0000000041                           Dr. Dexter Rolfson00000004470000000003     1563.4720210630",
					"0000000078                                    Wade Mraz00000008610000000003      224.9720210910",
					"0000000002                           Augustus Aufderhar00000000220000000000       190.820210530",
					"0000000025                             Frederica Cremin00000002760000000004      113.7520211103",
					"0000000069                             Dr. Tyree Rogahn00000007430000000002      1401.620210317",
					"0000000001                              Sammie Baumbach00000000070000000002       96.4720210528",
					"0000000077                         Mrs. Stephen Trantow00000008480000000004      1689.020210325",
					"0000000075                                  Bobbie Batz00000008120000000001      707.9620211103",
					"0000000070                              Palmer Prosacco00000007530000000003     1009.5420210308",
				}

				return bytes.NewBufferString(strings.Join(lines, "\n"))
			},
			inputHasHeader: true,
			wantResult: &model.LegacyImportResult{
				Users:    17,
				Orders:   20,
				Products: 21,
			},
			wantError: func() error {
				return nil
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				mockRepositoryOrder := new(mock_repository.MockRepositoryOrder)
				mockRepositoryOrder.On("LegacyBulkInsert").Return(nil)
				mockRepository.On("Order").Return(mockRepositoryOrder)

				mockCacheOrder := new(mock_cache.MockCacheOrder)
				mockCacheOrder.On("ClearAll").Return(nil)
				mockCache.On("Order").Return(mockCacheOrder)
			},
		},
		{
			name: "NotHasHeaderSuccess",
			inputFile: func() io.Reader {
				lines := []string{
					"0000000070                              Palmer Prosacco00000007530000000003     1836.7420210308",
					"0000000075                                  Bobbie Batz00000007980000000002     1578.5720211116",
					"0000000049                               Ken Wintheiser00000005230000000003      586.7420210903",
					"0000000014                                 Clelia Hills00000001460000000001      673.4920211125",
					"0000000057                          Elidia Gulgowski IV00000006200000000000     1417.2520210919",
					"0000000080                                 Tabitha Kuhn00000008770000000003      817.1320210612",
					"0000000023                                  Logan Lynch00000002530000000002      322.1220210523",
					"0000000015                                   Bonny Koss00000001530000000004        80.820210701",
					"0000000017                              Ethan Langworth00000001690000000000      865.1820210409",
					"0000000077                         Mrs. Stephen Trantow00000008440000000005     1288.7720211127",
					"0000000061                           Dimple Bergstrom I00000006710000000004       43.3620211104",
					"0000000077                         Mrs. Stephen Trantow00000008320000000006      961.3720210513",
					"0000000041                           Dr. Dexter Rolfson00000004470000000003     1563.4720210630",
					"0000000078                                    Wade Mraz00000008610000000003      224.9720210910",
					"0000000002                           Augustus Aufderhar00000000220000000000       190.820210530",
					"0000000025                             Frederica Cremin00000002760000000004      113.7520211103",
					"0000000069                             Dr. Tyree Rogahn00000007430000000002      1401.620210317",
					"0000000001                              Sammie Baumbach00000000070000000002       96.4720210528",
					"0000000077                         Mrs. Stephen Trantow00000008480000000004      1689.020210325",
					"0000000075                                  Bobbie Batz00000008120000000001      707.9620211103",
					"0000000070                              Palmer Prosacco00000007530000000003     1009.5420210308",
				}

				return bytes.NewBufferString(strings.Join(lines, "\n"))
			},
			inputHasHeader: false,
			wantResult: &model.LegacyImportResult{
				Users:    17,
				Orders:   20,
				Products: 21,
			},
			wantError: func() error {
				return nil
			},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				mockRepositoryOrder := new(mock_repository.MockRepositoryOrder)
				mockRepositoryOrder.On("LegacyBulkInsert").Return(nil)
				mockRepository.On("Order").Return(mockRepositoryOrder)

				mockCacheOrder := new(mock_cache.MockCacheOrder)
				mockCacheOrder.On("ClearAll").Return(nil)
				mockCache.On("Order").Return(mockCacheOrder)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := new(mock_repository.MockRepository)
			mockCache := new(mock_cache.MockCache)

			inputFile := tt.inputFile()
			wantError := tt.wantError()

			tt.mockOn(mockRepository, mockCache)

			usecaseOrder := NewOrder(mockRepository, mockCache)

			modelLegacyImportResult, err := usecaseOrder.LegacyImport(inputFile, tt.inputHasHeader)

			if !reflect.DeepEqual(err, wantError) {
				t.Errorf("Insert() got error = %v, want = %v.", err, wantError)
			}

			if !reflect.DeepEqual(modelLegacyImportResult, tt.wantResult) {
				t.Errorf("Insert() got result = %v, want = %v.", modelLegacyImportResult, tt.wantResult)
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
		name         string
		inputOrderId int64
		wantResult   *model.OrderDetails
		wantError    error
		mockOn       func(*mock_repository.MockRepository, *mock_cache.MockCache)
	}

	tests := []test{
		{
			name:         "RepositoryError",
			inputOrderId: 753,
			wantResult:   nil,
			wantError:    errors.New("Repository Error"),
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				mockRepositoryOrder := new(mock_repository.MockRepositoryOrder)
				mockRepositoryOrder.On("GetDetailsByOrderID").Return(nil, errors.New("Repository Error"))
				mockRepository.On("Order").Return(mockRepositoryOrder)

				mockCacheOrder := new(mock_cache.MockCacheOrder)
				mockCacheOrder.On("GetDetailsByOrderID").Return(nil, errors.New("Cache Error"))
				mockCache.On("Order").Return(mockCacheOrder)
			},
		},
		{
			name:         "CacheSuccess",
			inputOrderId: 753,
			wantResult:   &modelOrderDetails,
			wantError:    nil,
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				mockRepositoryOrder := new(mock_repository.MockRepositoryOrder)
				mockRepositoryOrder.On("GetDetailsByOrderID").Return(nil, nil)
				mockRepository.On("Order").Return(mockRepositoryOrder)

				mockCacheOrder := new(mock_cache.MockCacheOrder)
				mockCacheOrder.On("GetDetailsByOrderID").Return(&modelOrderDetails, nil)
				mockCache.On("Order").Return(mockCacheOrder)
			},
		},
		{
			name:         "RepositorySuccess",
			inputOrderId: 753,
			wantResult:   &modelOrderDetails,
			wantError:    nil,
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				mockRepositoryOrder := new(mock_repository.MockRepositoryOrder)
				mockRepositoryOrder.On("GetDetailsByOrderID").Return(&modelOrderDetails, nil)
				mockRepository.On("Order").Return(mockRepositoryOrder)

				mockCacheOrder := new(mock_cache.MockCacheOrder)
				mockCacheOrder.On("GetDetailsByOrderID").Return(nil, errors.New("Cache Error"))
				mockCacheOrder.On("SetDetailsByOrderID").Return(nil)
				mockCache.On("Order").Return(mockCacheOrder)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := new(mock_repository.MockRepository)
			mockCache := new(mock_cache.MockCache)

			tt.mockOn(mockRepository, mockCache)

			usecaseOrder := NewOrder(mockRepository, mockCache)

			modelLegacyImportResult, err := usecaseOrder.GetDetailsByOrderID(tt.inputOrderId)

			if !reflect.DeepEqual(err, tt.wantError) {
				t.Errorf("Insert() got error = %v, want = %v.", err, tt.wantError)
			}

			if !reflect.DeepEqual(modelLegacyImportResult, tt.wantResult) {
				t.Errorf("Insert() got result = %v, want = %v.", modelLegacyImportResult, tt.wantResult)
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
		name       string
		wantResult *model.OrdersDetails
		wantError  error
		mockOn     func(*mock_repository.MockRepository, *mock_cache.MockCache)
	}

	tests := []test{
		{
			name:       "RepositoryError",
			wantResult: nil,
			wantError:  errors.New("Repository Error"),
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				mockRepositoryOrder := new(mock_repository.MockRepositoryOrder)
				mockRepositoryOrder.On("ListDetails").Return(nil, errors.New("Repository Error"))
				mockRepository.On("Order").Return(mockRepositoryOrder)
			},
		},
		{
			name:       "Success",
			wantResult: &modelOrdersDetails,
			wantError:  nil,
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				mockRepositoryOrder := new(mock_repository.MockRepositoryOrder)
				mockRepositoryOrder.On("ListDetails").Return(&modelOrdersDetails, nil)
				mockRepository.On("Order").Return(mockRepositoryOrder)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := new(mock_repository.MockRepository)
			mockCache := new(mock_cache.MockCache)

			tt.mockOn(mockRepository, mockCache)

			usecaseOrder := NewOrder(mockRepository, mockCache)

			modelLegacyImportResult, err := usecaseOrder.ListDetails()

			if !reflect.DeepEqual(err, tt.wantError) {
				t.Errorf("Insert() got error = %v, want = %v.", err, tt.wantError)
			}

			if !reflect.DeepEqual(modelLegacyImportResult, tt.wantResult) {
				t.Errorf("Insert() got result = %v, want = %v.", modelLegacyImportResult, tt.wantResult)
			}
		})
	}
}

func TestOrderListDetailsByRangeBuyDate(t *testing.T) {
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
		name       string
		inputParam *model.OrderRangeBuyDate
		wantResult *model.OrdersDetails
		wantError  error
		mockOn     func(*mock_repository.MockRepository, *mock_cache.MockCache)
	}

	tests := []test{
		{
			name:       "ParamFromEmptyError",
			inputParam: &model.OrderRangeBuyDate{From: time.Time{}, To: OrderBuyDateMax},
			wantResult: nil,
			wantError:  ErrParamValidate{Message: OrderRangeBuyDateErrorMessageFromEmpty},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name:       "ParamFromBetweenError",
			inputParam: &model.OrderRangeBuyDate{From: OrderBuyDateMin.AddDate(0, 0, -1), To: OrderBuyDateMin},
			wantResult: nil,
			wantError:  ErrParamValidate{Message: OrderRangeBuyDateErrorMessageFromBetween},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name:       "ParamToEmptyError",
			inputParam: &model.OrderRangeBuyDate{From: OrderBuyDateMax, To: time.Time{}},
			wantResult: nil,
			wantError:  ErrParamValidate{Message: OrderRangeBuyDateErrorMessageToEmpty},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name:       "ParamToEmptyError",
			inputParam: &model.OrderRangeBuyDate{From: OrderBuyDateMax, To: OrderBuyDateMax.AddDate(0, 0, 1)},
			wantResult: nil,
			wantError:  ErrParamValidate{Message: OrderRangeBuyDateErrorMessageToBetween},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name:       "ParamToSmallerFromError",
			inputParam: &model.OrderRangeBuyDate{From: OrderBuyDateMax, To: OrderBuyDateMax.AddDate(0, 0, -1)},
			wantResult: nil,
			wantError:  ErrParamValidate{Message: OrderRangeBuyDateErrorMessageToSmallerFrom},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name:       "RangeError",
			inputParam: &model.OrderRangeBuyDate{From: OrderBuyDateMax.AddDate(0, 0, -32), To: OrderBuyDateMax},
			wantResult: nil,
			wantError:  ErrParamValidate{Message: OrderRangeBuyDateErrorMessageRangeError},
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				return
			},
		},
		{
			name:       "RepositoryError",
			inputParam: &model.OrderRangeBuyDate{From: OrderBuyDateMax, To: OrderBuyDateMax},
			wantResult: nil,
			wantError:  errors.New("Repository Error"),
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				mockRepositoryOrder := new(mock_repository.MockRepositoryOrder)
				mockRepositoryOrder.On("ListDetailsByRangeBuyDate").Return(nil, errors.New("Repository Error"))
				mockRepository.On("Order").Return(mockRepositoryOrder)
			},
		},
		{
			name:       "Success",
			inputParam: &model.OrderRangeBuyDate{From: OrderBuyDateMax, To: OrderBuyDateMax},
			wantResult: &modelOrdersDetails,
			wantError:  nil,
			mockOn: func(mockRepository *mock_repository.MockRepository, mockCache *mock_cache.MockCache) {
				mockRepositoryOrder := new(mock_repository.MockRepositoryOrder)
				mockRepositoryOrder.On("ListDetailsByRangeBuyDate").Return(&modelOrdersDetails, nil)
				mockRepository.On("Order").Return(mockRepositoryOrder)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := new(mock_repository.MockRepository)
			mockCache := new(mock_cache.MockCache)

			tt.mockOn(mockRepository, mockCache)

			usecaseOrder := NewOrder(mockRepository, mockCache)

			modelLegacyImportResult, err := usecaseOrder.ListDetailsByRangeBuyDate(tt.inputParam)

			if !reflect.DeepEqual(err, tt.wantError) {
				t.Errorf("Insert() got error = %v, want = %v.", err, tt.wantError)
			}

			if !reflect.DeepEqual(modelLegacyImportResult, tt.wantResult) {
				t.Errorf("Insert() got result = %v, want = %v.", modelLegacyImportResult, tt.wantResult)
			}
		})
	}
}
