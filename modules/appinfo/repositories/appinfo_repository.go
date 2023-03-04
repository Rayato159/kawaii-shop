package repositories

import (
	"fmt"
	"strings"

	"github.com/Rayato159/kawaii-shop/modules/appinfo"
	"github.com/jmoiron/sqlx"
)

type IAppinfoRepository interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
}

type appinfoRepository struct {
	db *sqlx.DB
}

func AppinfoRepository(db *sqlx.DB) IAppinfoRepository {
	return &appinfoRepository{
		db: db,
	}
}

func (r *appinfoRepository) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	query := `
	SELECT
		"id",
		"title"
	FROM "categories"`

	// Stack filter args
	filterValue := make([]any, 0)
	if req.Title != "" {
		query += `
		WHERE (LOWER("title") LIKE $1)`

		filterValue = append(filterValue, "%"+strings.ToLower(req.Title)+"%")
	}
	query += ";"

	category := make([]*appinfo.Category, 0)
	if err := r.db.Select(&category, query, filterValue...); err != nil {
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}
