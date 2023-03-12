package repositories

import (
	"github.com/Rayato159/kawaii-shop/modules/orders"
	"github.com/Rayato159/kawaii-shop/modules/orders/repositories/patterns"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
	FindOrders(req *orders.OrderFilter) ([]*orders.Order, int)
}

type ordersRepository struct {
	db *sqlx.DB
}

func OrdersRepository(db *sqlx.DB) IOrdersRepository {
	return &ordersRepository{
		db: db,
	}
}

func (r *ordersRepository) FindOrders(req *orders.OrderFilter) ([]*orders.Order, int) {
	builder := patterns.FindOrdersBuilder(r.db, req)
	engineer := patterns.FindOrdersEngineer(builder)

	return engineer.FindOrders(), engineer.CountOrders()
}
