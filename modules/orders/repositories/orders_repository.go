package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/Rayato159/kawaii-shop/modules/orders"
	"github.com/Rayato159/kawaii-shop/modules/orders/repositories/patterns"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
	FindOrders(req *orders.OrderFilter) ([]*orders.Order, error)
}

type ordersRepository struct {
	db *sqlx.DB
}

type iorderRows interface {
	scan() ([]*orders.Order, error)
}

func (o *orderRows) scan() ([]*orders.Order, error) {
	results := make([]*orders.Order, 0)

	for o.rows.Next() {
		// Init object
		data := make([]byte, 0)
		order := &orders.Order{
			TransterSlip: &orders.TransterSlip{},
			Products:     make([]*orders.ProductsOrder, 0),
		}
		// Scan
		if err := o.rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("scan order failed: %v", err)
		}
		// Unmarshal bytes to struct
		if err := json.Unmarshal(data, &order); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %v", err)
		}
		results = append(results, order)
	}
	return results, nil
}

type orderRows struct {
	rows *sqlx.Rows
}

func newOrderRows(rows *sqlx.Rows) iorderRows {
	return &orderRows{
		rows: rows,
	}
}

func OrdersRepository(db *sqlx.DB) IOrdersRepository {
	return &ordersRepository{
		db: db,
	}
}

func (r *ordersRepository) FindOrders(req *orders.OrderFilter) ([]*orders.Order, error) {
	builder := patterns.FindOrdersBuilder(r.db, req)
	engineer := patterns.FindOrdersEngineer(builder)
	engineer.FindOrders()
	return nil, nil
}
