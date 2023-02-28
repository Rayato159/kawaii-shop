package repositories

import "github.com/jmoiron/sqlx"

type IMiddlewareRepository interface{}

type middlewareRepository struct {
	Db *sqlx.DB
}

func MiddlewareRepository(db *sqlx.DB) IMiddlewareRepository {
	return &middlewareRepository{
		Db: db,
	}
}
