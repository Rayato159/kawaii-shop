package repositories

import (
	"fmt"

	"github.com/Rayato159/kawaii-shop/modules/appinfo"
	"github.com/jmoiron/sqlx"
)

type IAppinfoRepository interface {
	FindCategory() ([]*appinfo.Category, error)
}

type appinfoRepository struct {
	db *sqlx.DB
}

func AppinfoRepository(db *sqlx.DB) IAppinfoRepository {
	return &appinfoRepository{
		db: db,
	}
}

func (r *appinfoRepository) FindCategory() ([]*appinfo.Category, error) {
	query := `
	SELECT
		"id",
		"title"
	FROM "categories";`

	category := make([]*appinfo.Category, 0)
	if err := r.db.Select(&category, query); err != nil {
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}
