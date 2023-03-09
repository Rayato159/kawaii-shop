package repositories

import "github.com/jmoiron/sqlx"

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
