package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/Rayato159/kawaii-shop/modules/appinfo"
	"github.com/jmoiron/sqlx"
)

type IAppinfoRepository interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(categoryId int) error
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

func (r *appinfoRepository) InsertCategory(req []*appinfo.Category) error {
	query := `
	INSERT INTO "categories" (
		"title"
	)
	VALUES`

	tx, err := r.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}

	valuesStack := make([]any, 0)
	for i := range req {
		// Stack values
		valuesStack = append(valuesStack, req[i].Title)
		// Stack query
		if i != len(req)-1 {
			query += fmt.Sprintf(`
		($%d),`, i+1)
		} else {
			query += fmt.Sprintf(`
		($%d);`, i+1)
		}
	}

	if _, err := tx.ExecContext(context.Background(), query, valuesStack...); err != nil {
		tx.Rollback()
		return fmt.Errorf("insert many category failed: %v", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *appinfoRepository) DeleteCategory(categoryId int) error {
	query := `
	DELETE FROM "categories"
	WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, categoryId); err != nil {
		return fmt.Errorf("delete cateogry failed: %v", err)
	}
	return nil
}
