package patterns

import (
	"context"
	"fmt"

	"github.com/Rayato159/kawaii-shop/modules/products"
	"github.com/jmoiron/sqlx"
)

type IUpdateProductBuilder interface {
	initTransaction() error
	initQuery()
	updateTitleQuery()
	updateDescriptionQuery()
	updateCategoryQuery()
	searchOldImagesQuery()
	insertImagesQuery()
	deleteOldImagesQuery()
	closeQuery()
	updateProduct() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	commit() error
}

func (b *updateProductBuilder) getQueryFields() []string {
	return b.queryFields
}

func (b *updateProductBuilder) getValues() []any {
	return b.values
}

func (b *updateProductBuilder) getQuery() string {
	return b.query
}

func (b *updateProductBuilder) setQuery(query string) {
	b.query = query
}

func (b *updateProductBuilder) initQuery() {
	b.query += `
	UPDATE "products" SET`
}

func (b *updateProductBuilder) updateTitleQuery() {
	if b.req.Title != "" {
		b.values = append(b.values, b.req.Title)
		b.lastStackIndex = len(b.values)

		b.queryFields = append(b.queryFields, fmt.Sprintf(`
		"title" = $%d`, b.lastStackIndex))
	}
}

func (b *updateProductBuilder) updateDescriptionQuery() {
	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastStackIndex = len(b.values)

		b.queryFields = append(b.queryFields, fmt.Sprintf(`
		"description" = $%d`, b.lastStackIndex))
	}
}

func (b *updateProductBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}

func (b *updateProductBuilder) updateCategoryQuery() {
	b.query += `
	`
}

func (b *updateProductBuilder) searchOldImagesQuery() {
	b.query += `
	`
}

func (b *updateProductBuilder) insertImagesQuery() {
	b.query += `
	`
}

func (b *updateProductBuilder) deleteOldImagesQuery() {
	b.query += `
	`
}

func (b *updateProductBuilder) updateProduct() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update products failed: %v", err)
	}
	return nil
}

func (b *updateProductBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateProductBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		b.tx.Rollback()
		return err
	}
	return nil
}

type updateProductBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *products.Product
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdateProductBuilder(db *sqlx.DB, req *products.Product) IUpdateProductBuilder {
	return &updateProductBuilder{
		db:             db,
		req:            req,
		queryFields:    make([]string, 0),
		values:         make([]any, 0),
		lastStackIndex: 0,
	}
}

type updateProductEngineer struct {
	builder IUpdateProductBuilder
}

func UpdateProductEngineer(b IUpdateProductBuilder) *updateProductEngineer {
	return &updateProductEngineer{builder: b}
}

func (en *updateProductEngineer) sumQueryFieldsProducts() {
	en.builder.updateTitleQuery()
	en.builder.updateDescriptionQuery()

	fields := en.builder.getQueryFields()

	for i := range fields {
		query := en.builder.getQuery()
		if i != len(fields)-1 {
			en.builder.setQuery(query + fields[i] + ",")
		} else {
			en.builder.setQuery(query + fields[i])
		}
	}
}

func (en *updateProductEngineer) UpdateProduct() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFieldsProducts()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	// Commit
	if err := en.builder.updateProduct(); err != nil {
		return err
	}
	if err := en.builder.commit(); err != nil {
		return err
	}
	return nil
}
