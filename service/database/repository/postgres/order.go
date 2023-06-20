package repository

import (
	"database/sql"
	"fmt"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository"
	"github.com/lib/pq"
)

type PostgresOrder struct {
	Repository *Postgres
}

const (
	queryOrderDetails = `SELECT
			o.id, o.buy_date, o.total, o.user_id, u.name, op.product_id, op.product_value
		FROM
			orders o
		LEFT JOIN
			users u ON u.id = o.user_id
		LEFT JOIN
			orders_product op ON op.order_id = o.id 
		%s
		ORDER BY
			o.user_id, o.id`
)

func NewOrder(repository *Postgres) repository.Order {
	return &PostgresOrder{Repository: repository}
}

func (postgresOrder *PostgresOrder) ListDetails() (*model.OrdersDetails, error) {
	query := fmt.Sprintf(queryOrderDetails, "")

	rows, err := postgresOrder.Repository.Conn.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	modelOrdersDetails, err := postgresOrder.convertQueryResultToOrdersDetails(rows)

	// repository error not found
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			err = repository.ErrNotFound{Message: err.Error()}
		}
	} else if len(*modelOrdersDetails) == 0 {
		err = repository.ErrNotFound{Message: sql.ErrNoRows.Error()}
	}

	return modelOrdersDetails, err
}

func (postgresOrder *PostgresOrder) ListDetailsByRangeBuyDate(modelOrderRangeBuyDate *model.OrderRangeBuyDate) (*model.OrdersDetails, error) {
	query := fmt.Sprintf(queryOrderDetails, " WHERE o.buy_date BETWEEN $1 AND $2 ")

	rows, err := postgresOrder.Repository.Conn.Query(query, modelOrderRangeBuyDate.From, modelOrderRangeBuyDate.To)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	modelOrdersDetails, err := postgresOrder.convertQueryResultToOrdersDetails(rows)

	// repository error not found
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			err = repository.ErrNotFound{Message: err.Error()}
		}
	} else if len(*modelOrdersDetails) == 0 {
		err = repository.ErrNotFound{Message: sql.ErrNoRows.Error()}
	}

	return modelOrdersDetails, err
}

func (postgresOrder *PostgresOrder) GetDetailsByOrderID(orderID int64) (*model.OrderDetails, error) {
	query := fmt.Sprintf(queryOrderDetails, " WHERE o.id = $1 ")

	rows, err := postgresOrder.Repository.Conn.Query(query, orderID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	modelOrdersDetails, err := postgresOrder.convertQueryResultToOrdersDetails(rows)

	// repository error not found
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			err = repository.ErrNotFound{Message: err.Error()}
		}
	} else if len(*modelOrdersDetails) == 0 {
		err = repository.ErrNotFound{Message: sql.ErrNoRows.Error()}
	}

	if err != nil {
		return nil, err
	}

	return &(*modelOrdersDetails)[0], nil
}

func (postgresOrder *PostgresOrder) LegacyBulkInsert(modelUsers *model.Users, modelOrders *model.Orders, modelOrdersProducts *model.OrdersProducts) error {
	tx, err := postgresOrder.Repository.Conn.Begin()

	if err != nil {
		return err
	}

	err = postgresOrder.legacyClearAll(tx)

	if err == nil {
		err = postgresOrder.legacyUserBulkInsert(modelUsers, tx)
	}

	if err == nil {
		err = postgresOrder.legacyOrderBulkInsert(modelOrders, tx)
	}

	if err == nil {
		err = postgresOrder.legacyOrderProductBulkInsert(modelOrdersProducts, tx)
	}

	if err != nil {
		tx.Rollback()
	} else {
		err = tx.Commit()
	}

	// repository error duplicate key
	if errPQ, ok := err.(*pq.Error); ok {
		if errPQ.Code == "23505" {
			err = repository.ErrDuplicateKey{Message: errPQ.Detail}
		}
	}

	return err
}

func (postgresOrder *PostgresOrder) legacyClearAll(tx *sql.Tx) error {
	query := `TRUNCATE TABLE orders_product CASCADE;
		TRUNCATE TABLE orders CASCADE;
		TRUNCATE TABLE users CASCADE;`

	_, err := tx.Exec(query)

	return err
}

func (postgresOrder *PostgresOrder) legacyUserBulkInsert(modelUsers *model.Users, tx *sql.Tx) error {
	for _, modelUser := range *modelUsers {
		// userExists, err := postgresOrder.legacyUserCheckExistsByID(modelUser.ID)

		// if err != nil {
		// 	return err
		// }

		// if userExists {
		// 	continue
		// }

		err := postgresOrder.legacyUserInser(&modelUser, tx)

		if err != nil {
			return err
		}
	}

	return nil
}

func (postgresOrder *PostgresOrder) legacyUserCheckExistsByID(id int64) (exists bool, err error) {
	query :=
		`SELECT 
			id 
		FROM 
			users 
		WHERE 
			id = $1;`

	rows, err := postgresOrder.Repository.Conn.Query(query, id)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var userIDFound int64

		err = rows.Scan(&userIDFound)

		if err != nil {
			return
		}

		exists = (userIDFound == id)
	}

	return
}

func (*PostgresOrder) legacyUserInser(modelUser *model.User, tx *sql.Tx) error {
	query :=
		`INSERT INTO 
			users
			(id, name)
		VALUES
			($1, $2);`

	_, err := tx.Exec(
		query,
		modelUser.ID,
		modelUser.Name,
	)

	return err
}

func (postgresOrder *PostgresOrder) legacyOrderBulkInsert(modelOrders *model.Orders, tx *sql.Tx) (err error) {
	// var orderExists bool

	for _, modelOrder := range *modelOrders {
		// orderExists, err = postgresOrder.legacyOrderCheckExistsByID(modelOrder.ID)

		// if err != nil {
		// 	return err
		// }

		// if orderExists {
		// 	err = postgresOrder.legacyOrderUpdate(&modelOrder, tx)
		// } else {
		// 	err = postgresOrder.legacyOrderInsert(&modelOrder, tx)
		// }
		err = postgresOrder.legacyOrderInsert(&modelOrder, tx)

		if err != nil {
			return
		}
	}

	return
}

func (postgresOrder *PostgresOrder) legacyOrderCheckExistsByID(id int64) (exists bool, err error) {
	query :=
		`SELECT 
			id 
		FROM 
			orders 
		WHERE 
			id = $1;`

	rows, err := postgresOrder.Repository.Conn.Query(query, id)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var idFound int64

		err = rows.Scan(&idFound)

		if err != nil {
			return
		}

		exists = (idFound == id)
	}

	return
}

func (*PostgresOrder) legacyOrderInsert(modelOrder *model.Order, tx *sql.Tx) error {
	query :=
		`INSERT INTO 
			orders
			(id, user_id, buy_date, total)
		VALUES
			($1, $2, $3, $4);`

	_, err := tx.Exec(
		query,
		modelOrder.ID,
		modelOrder.UserID,
		modelOrder.BuyDate,
		modelOrder.Total,
	)

	return err
}

func (*PostgresOrder) legacyOrderUpdate(modelOrder *model.Order, tx *sql.Tx) error {
	query :=
		`UPDATE 
			orders 
		SET 
			total = total + $2
		WHERE
			id = $1;`

	_, err := tx.Exec(
		query,
		modelOrder.ID,
		modelOrder.Total,
	)

	return err
}

func (postgresOrder *PostgresOrder) legacyOrderProductBulkInsert(modelOrdersProducts *model.OrdersProducts, tx *sql.Tx) (err error) {
	for _, modelOrderProduct := range *modelOrdersProducts {
		err = postgresOrder.legacyOrderProductInsert(&modelOrderProduct, tx)

		if err != nil {
			return
		}
	}

	return
}

func (*PostgresOrder) legacyOrderProductInsert(modelOrderProduct *model.OrderProduct, tx *sql.Tx) error {
	query :=
		`INSERT INTO 
			orders_product
			(order_id, product_id, product_value)
		VALUES
			($1, $2, $3);`

	_, err := tx.Exec(
		query,
		modelOrderProduct.OrderID,
		modelOrderProduct.ProductID,
		modelOrderProduct.ProductValue,
	)

	return err
}

func (*PostgresOrder) convertQueryResultToOrdersDetails(rows *sql.Rows) (*model.OrdersDetails, error) {
	modelOrdersDetails := model.OrdersDetails{}
	userIndex := -1
	orderIndex := -1

	for rows.Next() {
		modelOrderUserProduct := model.OrderUserProduct{}

		err := rows.Scan(
			&modelOrderUserProduct.OrderID,
			&modelOrderUserProduct.OrderBuyDate,
			&modelOrderUserProduct.OrderTotal,
			&modelOrderUserProduct.UserID,
			&modelOrderUserProduct.UserName,
			&modelOrderUserProduct.ProductID,
			&modelOrderUserProduct.ProductValue,
		)

		if err != nil {
			return nil, err
		}

		if userIndex < 0 || modelOrdersDetails[userIndex].UserID != modelOrderUserProduct.UserID {
			modelOrdersDetails = append(modelOrdersDetails, model.OrderDetails{
				UserID:   modelOrderUserProduct.UserID,
				UserName: modelOrderUserProduct.UserName,
			})

			orderIndex = -1
			userIndex++
		}

		if orderIndex < 0 || modelOrdersDetails[userIndex].Orders[orderIndex].OrderID != modelOrderUserProduct.OrderID {
			modelOrdersDetails[userIndex].Orders = append(modelOrdersDetails[userIndex].Orders, model.OrderDetailsOrder{
				OrderID: modelOrderUserProduct.OrderID,
				BuyDate: modelOrderUserProduct.OrderBuyDate.Format("2006-01-02"),
				Total:   modelOrderUserProduct.OrderTotal,
			})

			orderIndex++
		}

		modelOrdersDetails[userIndex].Orders[orderIndex].Products = append(modelOrdersDetails[userIndex].Orders[orderIndex].Products, model.OrderDetailsProduct{
			ID:    modelOrderUserProduct.ProductID,
			Value: modelOrderUserProduct.ProductValue,
		})
	}

	return &modelOrdersDetails, nil
}
