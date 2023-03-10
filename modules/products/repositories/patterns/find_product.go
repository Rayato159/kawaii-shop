package patterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Rayato159/kawaii-shop/modules/products"
	"github.com/Rayato159/kawaii-shop/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type IFindProductBuidler interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	closeJsonQuery()
	sort()
	paginate()
	resetQuery()
	Result() []*products.Product
	Count() int
	PrintQuery()
}

func (b *findProductBuilder) openJsonQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM (`
}

func (b *findProductBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}

func (b *findProductBuilder) initQuery() {
	b.query += `
		SELECT
			"p"."id",
			"p"."title",
			"p"."description",
			"p"."price",
			(
				SELECT
					to_json("ct")
				FROM (
					SELECT
						"c"."id",
						"c"."title"
					FROM "categories" "c"
						LEFT JOIN "products_categories" "pc" ON "pc"."category_id" = "c"."id"
					WHERE "pc"."product_id" = "p"."id"
				) AS "ct"
			) AS "category",
			"p"."created_at",
			"p"."updated_at",
			(
				SELECT
					array_to_json(array_agg("it"))
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "images" "i"
					WHERE "i"."product_id" = "p"."id"
				) AS "it"
			) AS "images"
		FROM "products" "p"
		WHERE 1 = 1`
}

func (b *findProductBuilder) countQuery() {
	b.query += `
		SELECT
			COUNT(*) AS "count"
		FROM "products" "p"
		WHERE 1 = 1`
}

func (b *findProductBuilder) whereQuery() {
	// Where logic
	var queryWhere string
	queryWhereStack := make([]string, 0)
	if b.req.Id != "" {
		b.values = append(b.values, b.req.Id)

		queryWhereStack = append(queryWhereStack, `
		AND "p"."id" = ?`)
	}
	if b.req.Search != "" {
		b.values = append(b.values, "%"+b.req.Search+"%", "%"+b.req.Search+"%")

		queryWhereStack = append(queryWhereStack, `
		AND (LOWER("p"."title") LIKE ? OR LOWER("p"."description") LIKE ?)`)
	}

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1)
		} else {
			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1)
			queryWhere = strings.Replace(queryWhere, "?", "$"+strconv.Itoa(i+2), 1)
		}
	}
	b.lastStackIndex = len(b.values)

	// Sum total query
	b.query += queryWhere
}

func (b *findProductBuilder) closeJsonQuery() {
	b.query += `
	) AS "t";`
}

func (b *findProductBuilder) sort() {
	orderByMap := map[string]string{
		"id":    "\"p\".\"id\"",
		"title": "\"p\".\"title\"",
		"price": "\"p\".\"price\"",
	}
	if orderByMap[b.req.OrderBy] == "" {
		b.req.OrderBy = orderByMap["title"]
	} else {
		b.req.OrderBy = orderByMap[b.req.OrderBy]
	}

	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[b.req.Sort] == "" {
		b.req.Sort = sortMap["asc"]
	} else {
		b.req.Sort = sortMap[strings.ToUpper(b.req.Sort)]
	}

	b.values = append(b.values, b.req.OrderBy)
	b.query += fmt.Sprintf(`
		ORDER BY $%d %s`, b.lastStackIndex+1, b.req.Sort)
	b.lastStackIndex = len(b.values)
}

func (b *findProductBuilder) paginate() {
	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)

	b.query += fmt.Sprintf(` OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)
}

func (b *findProductBuilder) Result() []*products.Product {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	products := make([]*products.Product, 0)
	bytes := make([]byte, 0)

	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find products failed: %v\n", err)
		return products
	}

	if err := json.Unmarshal(bytes, &products); err != nil {
		log.Printf("unmarsal products failed: %v\n", err)
		return products
	}
	b.resetQuery()
	return products
}

func (b *findProductBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count products failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}

func (b *findProductBuilder) PrintQuery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}

type findProductBuilder struct {
	db             *sqlx.DB
	req            *products.ProductFilter
	query          string
	lastStackIndex int
	values         []any
}

func FindProductBuilder(db *sqlx.DB, req *products.ProductFilter) IFindProductBuidler {
	return &findProductBuilder{
		db:  db,
		req: req,
	}
}

type findProductEngineer struct {
	builder IFindProductBuidler
}

func FindProductEngineer(b IFindProductBuidler) *findProductEngineer {
	return &findProductEngineer{builder: b}
}

func (en *findProductEngineer) FindProduct() IFindProductBuidler {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findProductEngineer) CountProduct() IFindProductBuidler {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}
