package patterns

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/Rayato159/kawaii-shop/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IFindOrdersBuilder interface {
	initQuery()
	initCountQuery()
	productQuery()
	buildWhereSearch()
	buildWhereStatus()
	buildWhereDate()
	buildSort()
	buildPaginate()
	finalQuery()
	closeQuery()
	getQuery() string
	setQuery(query string)
	getValues() []any
	setValues(data []any)
	setLastIndex(n int)
	getDb() *sqlx.DB
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

func (b *findOrdersBuilder) getDb() *sqlx.DB {
	return b.db
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

func (b *findOrdersBuilder) setValues(data []any) {
	b.values = data
}

func (b *findOrdersBuilder) setLastIndex(n int) {
	b.lastIndex = n
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

func (b *findOrdersBuilder) initCountQuery() {
	b.query += `
	SELECT
		COUNT(*) AS "count"
	FROM "orders" "o"
	WHERE 1 = 1`
}

func (b *findOrdersBuilder) productQuery() {
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
					WHERE "spo"."order_id" = "o"."id"
				) AS "pt"
			) AS "products",`
}

func (b *findOrdersBuilder) buildWhereStatus() {
	if b.req.Search != "" {
		b.values = append(
			b.values,
			strings.ToLower(b.req.Status),
		)

		query := `
		AND "o"."status" = $?`

		temp := b.getQuery()
		query = strings.Replace(query, "?", strconv.Itoa(b.lastIndex+1), 1)
		temp += query
		b.setQuery(temp)

		b.lastIndex = len(b.values)
	}
}

func (b *findOrdersBuilder) buildWhereSearch() {
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		query := `
		AND (
			LOWER("o"."user_id") LIKE $? OR
			LOWER("o"."address") LIKE $? OR
			LOWER("o"."contact") LIKE $?
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

		b.lastIndex = len(b.values)
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

func (b *findOrdersBuilder) buildSort() {
	sortMap := map[string]string{
		"id":         `"o"."id"`,
		"created_at": `"o"."created_at"`,
	}
	if sortMap[b.req.OrderBy] == "" {
		b.req.OrderBy = `"o"."id"`
	}
	if b.req.Sort == "" {
		b.req.Sort = "DESC"
	}
	b.req.OrderBy = sortMap[b.req.OrderBy]

	b.values = append(
		b.values,
		b.req.OrderBy,
	)

	b.query += fmt.Sprintf(`
	ORDER BY $%d %s`, b.lastIndex+1, strings.ToUpper(b.req.Sort))

	b.lastIndex = len(b.values)
}

func (b *findOrdersBuilder) buildPaginate() {
	b.values = append(
		b.values,
		b.req.PaginateReq.Limit,
		math.Ceil(float64((b.req.PaginateReq.Page-1))*float64(b.req.PaginateReq.Limit)),
	)

	b.query += fmt.Sprintf(`
	LIMIT $%d OFFSET $%d`, b.lastIndex+1, b.lastIndex+2)

	b.lastIndex = len(b.values)
}

func (b *findOrdersBuilder) closeQuery() {
	b.query += `
) AS "t"`
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

func (en *findOrdersEngineer) FindOrders() []*orders.Order {
	en.builder.initQuery()
	en.builder.productQuery()
	en.builder.finalQuery()
	en.builder.buildWhereStatus()
	en.builder.buildWhereSearch()
	en.builder.buildWhereDate()
	en.builder.buildSort()
	en.builder.buildPaginate()
	en.builder.closeQuery()

	rows, err := en.builder.getDb().Queryx(en.builder.getQuery(), en.builder.getValues()...)
	if err != nil {
		log.Printf("orders query rows failed: %v", err)
		return make([]*orders.Order, 0)
	}

	results, err := newOrderRows(rows).scan()
	if err != nil {
		log.Println(err)
		return make([]*orders.Order, 0)
	}

	en.builder.setQuery("")
	en.builder.setValues(make([]any, 0))
	en.builder.setLastIndex(0)
	return results
}

func (en *findOrdersEngineer) CountOrders() int {
	en.builder.initCountQuery()
	en.builder.buildWhereStatus()
	en.builder.buildWhereSearch()
	en.builder.buildWhereDate()

	var count int
	if err := en.builder.getDb().Get(&count, en.builder.getQuery(), en.builder.getValues()...); err != nil {
		log.Printf("count orders failed: %v\n", err)
		return 0
	}

	en.builder.setQuery("")
	en.builder.setValues(make([]any, 0))
	en.builder.setLastIndex(0)
	return count
}
