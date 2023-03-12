package patterns

import (
	"context"
	"fmt"

	"github.com/Rayato159/kawaii-shop/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IInsertOrder interface {
	initTransaction() error
	insertOrder() error
	insertProductsOrder() error
	commit() error
	getOrderId() string
}

type insertOrderBuilder struct {
	req *orders.Order
	db  *sqlx.DB
	tx  *sqlx.Tx
}

type insertOrderEngineer struct {
	builder IInsertOrder
}

func (b *insertOrderBuilder) getOrderId() string {
	return b.req.Id
}

func (b *insertOrderBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertOrderBuilder) insertOrder() error {
	query := `
	INSERT INTO "orders" (
		"user_id",
		"contact",
		"address",
		"transfer_slip",
		"status"
	)
	VALUES
	(
		$1,
		$2,
		$3,
		$4,
		$5
	)
		RETURNING "id";`

	if err := b.tx.QueryRowxContext(
		context.Background(),
		query,
		b.req.UserId,
		b.req.Contact,
		b.req.Address,
		b.req.TransterSlip,
		b.req.Status,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert order failed: %v", err)
	}
	return nil
}

func (b *insertOrderBuilder) insertProductsOrder() error {
	query := `
	INSERT INTO "products_orders" (
		"order_id",
		"qty",
		"product"
	)
	VALUES`

	valuesStack := make([]any, 0)

	lastIndex := 0
	for i := range b.req.Products {
		valuesStack = append(
			valuesStack,
			b.req.Id,
			b.req.Products[i].Qty,
			b.req.Products[i].Product,
		)

		if i != len(b.req.Products)-1 {
			query += fmt.Sprintf(`
		($%d, $%d, $%d),`, lastIndex+1, lastIndex+2, lastIndex+3)
		} else {
			query += fmt.Sprintf(`
		($%d, $%d, $%d);`, lastIndex+1, lastIndex+2, lastIndex+3)
		}

		lastIndex += 3
	}

	if _, err := b.tx.ExecContext(context.Background(), query, valuesStack...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert products_orders failed: %v", err)
	}
	return nil
}

func (b *insertOrderBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func InsertOrderBuilder(db *sqlx.DB, req *orders.Order) IInsertOrder {
	return &insertOrderBuilder{
		db:  db,
		req: req,
	}
}

func InsertOrderEngineer(b IInsertOrder) *insertOrderEngineer {
	return &insertOrderEngineer{
		builder: b,
	}
}

func (en *insertOrderEngineer) InsertOrder() (string, error) {
	en.builder.initTransaction()
	en.builder.insertOrder()
	en.builder.insertProductsOrder()
	en.builder.commit()

	return en.builder.getOrderId(), nil
}
