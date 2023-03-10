package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Rayato159/kawaii-shop/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
	FindOrders(req *orders.OrderFilter) ([]*orders.Order, error)
}

type ordersRepository struct {
	db *sqlx.DB
}

type iorderRows interface {
	scan() ([]*orders.Order, error)
}

func (o *orderRows) scan() ([]*orders.Order, error) {
	results := make([]*orders.Order, 0)

	for o.rows.Next() {
		// Init object
		data := make([]byte, 0)
		order := &orders.Order{
			TransterSlip: &orders.TransterSlip{},
			Products:     make([]*orders.ProductsOrder, 0),
		}
		// Scan
		if err := o.rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("scan order failed: %v", err)
		}
		// Unmarshal bytes to struct
		if err := json.Unmarshal(data, &order); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %v", err)
		}
		results = append(results, order)
	}
	return results, nil
}

type orderRows struct {
	rows *sqlx.Rows
}

func newOrderRows(rows *sqlx.Rows) iorderRows {
	return &orderRows{
		rows: rows,
	}
}

func OrdersRepository(db *sqlx.DB) IOrdersRepository {
	return &ordersRepository{
		db: db,
	}
}

func (r *ordersRepository) FindOrders(req *orders.OrderFilter) ([]*orders.Order, error) {
	queryWhereStack := make([]string, 0)
	values := make([]any, 0)

	if req.Search != "" {
		values = append(
			values,
			"%"+strings.ToLower(req.Search)+"%",
		)
		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		AND (
				LOWER("o"."user_id") LIKE $? AND
				LOWER("o"."address") LIKE $? AND
				LOWER("o"."contract") LIKE $?
			)`))
	}

	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"o"."id",
			"o"."user_id",
			(
				SELECT
					array_to_json(array_agg("pt"))
				FROM (
					SELECT
						"spo"."id",
						"spo"."qty",
						"spo"."product"
					FROM "products_orders" "spo"
					WHERE "spo"."order_id" = "o"."id"
				) AS "pt"
			) AS "products",
			"o"."transfer_slip",
			"o"."contact",
			"o"."address",
			"o"."status",
			(
				SELECT
					SUM(COALESCE(("po"."product"->>'price')::FLOAT*("po"."qty")::FLOAT, 0))
				FROM "products_orders" "po"
				WHERE "po"."order_id" = "o"."id"
			) AS "total_paid",
			"o"."created_at",
			"o"."updated_at"
		FROM "orders" "o"
		WHERE 1 = 1`

	for i := range queryWhereStack {
		query += queryWhereStack[i]
	}
	query += `
	) AS "t"`

	rows, err := r.db.QueryxContext(
		context.Background(),
		query,
		values...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders, err := newOrderRows(rows).scan()
	if err != nil {
		return nil, err
	}
	return orders, nil
}
