package repositories

import "github.com/jmoiron/sqlx"

type IOauthRepository interface{}

type oauthRepository struct {
	Db *sqlx.DB
}

func OauthRepository(db *sqlx.DB) IOauthRepository {
	return &oauthRepository{Db: db}
}
