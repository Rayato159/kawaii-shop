package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Rayato159/kawaii-shop/modules/orders"
	"github.com/Rayato159/kawaii-shop/modules/orders/repositories/patterns"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
	FindOrder(req *orders.OrderFilter) ([]*orders.Order, int)
	FindOneOrder(orderId string) (*orders.Order, error)
	InsertOrder(req *orders.Order) (string, error)
	UpdateOrder(req *orders.UpdateOrderReq) error
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

func (r *ordersRepository) InsertOrder(req *orders.Order) (string, error) {
	builder := patterns.InsertOrderBuilder(r.db, req)
	orderId, err := patterns.InsertOrderEngineer(builder).InsertOrder()
	if err != nil {
		return "", err
	}
	return orderId, nil
}

func (r *ordersRepository) UpdateOrder(req *orders.UpdateOrderReq) error {
	query := `
	UPDATE "orders" SET`

	queryWhereStack := make([]string, 0)
	valueStack := make([]any, 0)
	lastIndex := 1

	if req.Status != "" {
		valueStack = append(valueStack, req.Status)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"status" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.TransterSlip != nil {
		valueStack = append(valueStack, req.TransterSlip)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"transfer_slip" = $%d?`, lastIndex))

		lastIndex++
	}

	valueStack = append(valueStack, req.OrderId)

	queryClose := fmt.Sprintf(`
	WHERE "id" = $%d;`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1)
		}
	}
	query += queryClose

	if _, err := r.db.ExecContext(context.Background(), query, valueStack...); err != nil {
		return fmt.Errorf("update order failed: %v", err)
	}
	return nil
}
