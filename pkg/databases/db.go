package databases

import (
	"log"

	"github.com/Rayato159/kawaii-shop/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	// Connect
	db, err := sqlx.Connect("pgx", cfg.Url())
	defer db.Close()
	if err != nil {
		db.Close()
		log.Fatalf("connect to database failed: %v\n", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())
	return db
}
