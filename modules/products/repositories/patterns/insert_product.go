package patterns

import (
	"context"
	"fmt"
	"time"

	"github.com/Rayato159/kawaii-shop/modules/products"
	"github.com/jmoiron/sqlx"
)

type IInsertProductBuilder interface {
	initTransaction() error
	insertProduct() error
	insertCategory() error
	insertAttachment() error
	commit() error
	getProductId() string
}

func (b *insertProductBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertProductBuilder) insertProduct() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO "products" (
		"title",
		"description"
	)
	VALUES ($1, $2)
		RETURNING "id";`

	if err := b.tx.QueryRowxContext(
		ctx,
		query,
		b.req.Title,
		b.req.Description,
	).Scan(&b.productId); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert product failed: %v", err)
	}
	return nil
}

func (b *insertProductBuilder) insertAttachment() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO "images" (
		"filename",
		"url",
		"product_id"
	)
	VALUES`

	valuesStack := make([]any, 0)
	var index int
	for i := range b.req.Images {
		valuesStack = append(
			valuesStack,
			b.req.Images[i].FileName,
			b.req.Images[i].Url,
			b.productId,
		)

		if i != len(b.req.Images)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
		index += 3
	}

	fmt.Println(query)

	if _, err := b.tx.ExecContext(
		ctx,
		query,
		valuesStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert attachments failed: %v", err)
	}
	return nil
}

func (b *insertProductBuilder) insertCategory() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO "products_categories" (
		"product_id",
		"category_id"
	)
	VALUES ($1, $2);`

	if _, err := b.tx.ExecContext(
		ctx,
		query,
		b.productId,
		b.req.Category.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert products_categories failed: %v", err)
	}
	return nil
}

func (b *insertProductBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		b.tx.Rollback()
		return err
	}
	return nil
}

func (b *insertProductBuilder) getProductId() string {
	return b.productId
}

type insertProductBuilder struct {
	db        *sqlx.DB
	tx        *sqlx.Tx
	req       *products.Product
	productId string
}

func InsertProductBuilder(db *sqlx.DB, req *products.Product) IInsertProductBuilder {
	return &insertProductBuilder{
		db:  db,
		req: req,
	}
}

type insertProductEngineer struct {
	builder IInsertProductBuilder
}

func InsertProductEngineer(b IInsertProductBuilder) *insertProductEngineer {
	return &insertProductEngineer{
		builder: b,
	}
}

func (en *insertProductEngineer) InsertProduct() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertProduct(); err != nil {
		return "", err
	}
	if err := en.builder.insertCategory(); err != nil {
		return "", err
	}
	if err := en.builder.insertAttachment(); err != nil {
		return "", err
	}
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getProductId(), nil
}
