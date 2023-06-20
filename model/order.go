package model

import "time"

type Legacy struct {
	UserID       int64
	UserName     string
	OrderID      int64
	ProductID    int64
	ProductValue float64
	BuyDate      string
	ImportedAt   time.Time
}

type LegacyRecord struct {
	UserID       string
	UserName     string
	OrderID      string
	ProductID    string
	ProductValue string
	BuyDate      string
}

type LegacyRecordError struct {
	Line    int64  `json:"line"`
	Message string `json:"message"`
}

type LegacyRecordsError []LegacyRecordError

type LegacyImportResult struct {
	// Quantidade de usuários importados
	Users int `json:"users" validate:"required"`
	// Quantidade de pedidos importados
	Orders int `json:"orders" validate:"required"`
	// Quantidade de produtos importados
	Products int `json:"products" validate:"required"`
}

type User struct {
	ID   int64
	Name string
}

type Users []User

type Order struct {
	ID      int64
	UserID  int64
	BuyDate string
	Total   float64
}

type Orders []Order

type OrderProduct struct {
	OrderID      int64
	ProductID    int64
	ProductValue float64
}

type OrdersProducts []OrderProduct

type OrderUserProduct struct {
	OrderID      int64
	OrderBuyDate time.Time
	OrderTotal   float64
	UserID       int64
	UserName     string
	ProductID    int64
	ProductValue float64
}

type OrderDetails struct {
	// ID do Usuário
	UserID int64 `json:"user_id" validate:"required" example:"1"`
	// Nome do Usuário
	UserName string `json:"name" validate:"required" example:"Joao"`
	// Lista de Pedidos
	Orders []OrderDetailsOrder `json:"orders" validate:"required"`
}

type OrderDetailsOrder struct {
	// ID do Pedido
	OrderID int64 `json:"order_id" validate:"required" example:"1"`
	// Data da Compra
	BuyDate string `json:"date" validate:"required" example:"2019-08-24" format:"date"`
	// Valor Total do Pedido
	Total float64 `json:"total" validate:"required" example:"23.45" format:"float"`
	// Lista de Produtos
	Products []OrderDetailsProduct `json:"products" validate:"required"`
}

type OrderDetailsProduct struct {
	// ID do Produto
	ID int64 `json:"product_id" validate:"required" example:"1"`
	// Valor do Produto
	Value float64 `json:"value" validate:"required" example:"23.45" format:"float"`
}

type OrdersDetails []OrderDetails

type OrderRangeBuyDate struct {
	From time.Time
	To   time.Time
}
