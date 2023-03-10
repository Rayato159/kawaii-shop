package patterns

import (
	"github.com/Rayato159/kawaii-shop/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IFindOrdersBuilder interface {
	initQuery()
	productQueryTop()
	productQueryButtom()
	injectWhereProduct()
	injectWhereSearch()
	finalQuery()
	getQuery() string
	setQuery(query string)
	getQueryWhereStack() []string
	getValues() []any
}

type findOrdersBuilder struct {
	db              *sqlx.DB
	req             *orders.OrderFilter
	query           string
	queryWhereStack []string
	values          []any
	lastIndex       int
}

type findOrdersEngineer struct {
	builder IFindOrdersBuilder
}

func (b *findOrdersBuilder) getQuery() string {
	return b.query
}

func (b *findOrdersBuilder) setQuery(query string) {
	b.query = query
}

func (b *findOrdersBuilder) getQueryWhereStack() []string {
	return b.queryWhereStack
}

func (b *findOrdersBuilder) getValues() []any {
	return b.values
}

func (b *findOrdersBuilder) initQuery() {
	b.query += `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"o"."id",
			"o"."user_id",`
}

func (b *findOrdersBuilder) productQueryTop() {
	b.query += `
	(
		SELECT
			array_to_json(array_agg("pt"))
		FROM (
			SELECT
				"spo"."id",
				"spo"."qty",
				"spo"."product"
			FROM "products_orders" "spo"
			WHERE "spo"."order_id" = "o"."id"`
}

func (b *findOrdersBuilder) injectWhereProduct() {
	b.query += `
	`
}

func (b *findOrdersBuilder) productQueryButtom() {
	b.query += `
		) AS "pt"
	) AS "products"`
}

func (b *findOrdersBuilder) injectWhereSearch() {
	b.query += `
	`
}

func (b *findOrdersBuilder) finalQuery() {
	b.query += `
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
}

func FindOrdersBuilder(db *sqlx.DB, req *orders.OrderFilter) IFindOrdersBuilder {
	return &findOrdersBuilder{
		db:              db,
		req:             req,
		queryWhereStack: make([]string, 0),
		values:          make([]any, 0),
		lastIndex:       0,
	}
}

func FindOrdersEngineer(b IFindOrdersBuilder) *findOrdersEngineer {
	return &findOrdersEngineer{
		builder: b,
	}
}

func (en *findOrdersEngineer) FindOrders() ([]*orders.Order, error) {
	return nil, nil
}
