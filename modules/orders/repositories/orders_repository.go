package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/Rayato159/kawaii-shop/modules/orders"
	"github.com/Rayato159/kawaii-shop/modules/orders/repositories/patterns"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
	FindOrder(req *orders.OrderFilter) ([]*orders.Order, int)
	FindOneOrder(orderId string) (*orders.Order, error)
}

type ordersRepository struct {
	db *sqlx.DB
}

func OrdersRepository(db *sqlx.DB) IOrdersRepository {
	return &ordersRepository{
		db: db,
	}
}

func (r *ordersRepository) FindOrder(req *orders.OrderFilter) ([]*orders.Order, int) {
	builder := patterns.FindOrdersBuilder(r.db, req)
	engineer := patterns.FindOrdersEngineer(builder)

	return engineer.FindOrders(), engineer.CountOrders()
}

func (r *ordersRepository) FindOneOrder(orderId string) (*orders.Order, error) {
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
		WHERE "o"."id" = $1
	) AS "t";`

	order := &orders.Order{
		Products:     make([]*orders.ProductsOrder, 0),
		TransterSlip: &orders.TransterSlip{},
	}

	raw := make([]byte, 0)
	if err := r.db.Get(&raw, query, orderId); err != nil {
		return nil, fmt.Errorf("get order failed: %v", err)
	}

	if err := json.Unmarshal(raw, &order); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %v", err)
	}
	return order, nil
}
