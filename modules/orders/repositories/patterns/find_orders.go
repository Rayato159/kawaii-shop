package patterns

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Rayato159/kawaii-shop/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IFindOrdersBuilder interface {
	initQuery()
	productQueryTop()
	productQueryButtom()
	buildWhereProduct()
	buildWhereSearch()
	buildWhereDate()
	finalQuery()
	closeQuery()
	getQuery() string
	setQuery(query string)
	getValues() []any
}

type findOrdersBuilder struct {
	db        *sqlx.DB
	req       *orders.OrderFilter
	query     string
	values    []any
	lastIndex int
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

func (b *findOrdersBuilder) buildWhereProduct() {
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+b.req.Search+"%",
			"%"+b.req.Search+"%",
			"%"+b.req.Search+"%",
		)

		query := `
					AND (
						LOWER("spo"."product"->>'id') LIKE $? OR
						LOWER("spo"."product"->>'title') LIKE $? OR
						LOWER("spo"."product"->>'description') LIKE $?
					)`

		temp := b.getQuery()
		query = strings.Replace(query, "?", strconv.Itoa(b.lastIndex+1), 1)
		query = strings.Replace(query, "?", strconv.Itoa(b.lastIndex+2), 1)
		query = strings.Replace(query, "?", strconv.Itoa(b.lastIndex+3), 1)
		temp += query
		b.setQuery(temp)

		b.lastIndex = len(b.values)
	}
}

func (b *findOrdersBuilder) productQueryButtom() {
	b.query += `
				) AS "pt"
			) AS "products",`
}

func (b *findOrdersBuilder) buildWhereSearch() {
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+b.req.Search+"%",
			"%"+b.req.Search+"%",
			"%"+b.req.Search+"%",
		)

		query := `
		AND (
			LOWER("o"."user_id") LIKE $? AND
			LOWER("o"."address") LIKE $? AND
			LOWER("o"."contract") LIKE $?
		)`

		temp := b.getQuery()
		query = strings.Replace(query, "?", strconv.Itoa(b.lastIndex+1), 1)
		query = strings.Replace(query, "?", strconv.Itoa(b.lastIndex+2), 1)
		query = strings.Replace(query, "?", strconv.Itoa(b.lastIndex+3), 1)
		temp += query
		b.setQuery(temp)

		b.lastIndex = len(b.values)
	}
}

func (b *findOrdersBuilder) buildWhereDate() {
	if b.req.StartDate != "" && b.req.EndDate != "" {
		b.values = append(
			b.values,
			"%"+b.req.StartDate+"%",
			"%"+b.req.EndDate+"%",
		)

		b.query += fmt.Sprintf(`
		AND "o"."created_at" BETWEEN $%d AND $%d`, b.lastIndex+1, b.lastIndex+2)
	}
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

func (b *findOrdersBuilder) closeQuery() {
	b.query += `
) AS "t"`
}

func FindOrdersBuilder(db *sqlx.DB, req *orders.OrderFilter) IFindOrdersBuilder {
	return &findOrdersBuilder{
		db:     db,
		req:    req,
		values: make([]any, 0),
	}
}

func FindOrdersEngineer(b IFindOrdersBuilder) *findOrdersEngineer {
	return &findOrdersEngineer{
		builder: b,
	}
}

func (en *findOrdersEngineer) FindOrders() ([]*orders.Order, error) {
	en.builder.initQuery()
	en.builder.productQueryTop()
	en.builder.buildWhereProduct()
	en.builder.productQueryButtom()
	en.builder.finalQuery()
	en.builder.buildWhereSearch()
	en.builder.buildWhereDate()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())
	return nil, nil
}
