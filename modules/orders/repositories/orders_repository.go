package repositories

import (
	"github.com/Rayato159/kawaii-shop/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
}

type ordersRepository struct {
	db *sqlx.DB
}

func OrdersRepository(db *sqlx.DB) IOrdersRepository {
	return &ordersRepository{
		db: db,
	}
}

func (r *ordersRepository) FindOrders(req *orders.OrderFilter) ([]*orders.Order, error) {
	return nil, nil
}
