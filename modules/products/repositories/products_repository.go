package repositories

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/products"

	"github.com/jmoiron/sqlx"
)

type IProductsRepository interface {
	FindProduct(req *products.ProductFilter) []*products.Product
	FindOneProduct(productId string) (*products.Product, error)
}

type productsRepository struct {
	db *sqlx.DB
}

func ProductRepository(db *sqlx.DB) IProductsRepository {
	return &productsRepository{
		db: db,
	}
}

// Need to refactor in patterns: builder (paginate needed too)
func (r *productsRepository) FindProduct(req *products.ProductFilter) []*products.Product {
	// Total query
	var query string

	// Init
	queryInit := `
	SELECT
		array_to_json(array_agg("t"))
	FROM (
		SELECT
			"p"."id",
			"p"."title",
			"p"."description",
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

	// Where logic
	var queryWhere string
	valuesStack := make([]any, 0)
	queryWhereStack := make([]string, 0)
	if req.Id != "" {
		valuesStack = append(valuesStack, req.Id)

		queryWhereStack = append(queryWhereStack, `
		AND "p"."id" = ?`)
	}
	if req.Search != "" {
		valuesStack = append(valuesStack, "%"+req.Search+"%", "%"+req.Search+"%")

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

	// Close
	queryClose := `
	) AS "t";`

	// Concat query
	query = queryInit + queryWhere + queryClose
	fmt.Println(query)

	productsBytes := make([]byte, 0)
	products := make([]*products.Product, 0)

	// Query in bytes
	if err := r.db.Get(&productsBytes, query, valuesStack...); err != nil {
		log.Printf("get products failed: %v", err)
		return products
	}

	// Parse bytes to json object
	if err := json.Unmarshal(productsBytes, &products); err != nil {
		log.Printf("unmarshal products failed: %v", err)
		return products
	}

	return products
}

// Need to refactor in patterns: builder
func (r *productsRepository) FindOneProduct(productId string) (*products.Product, error) {
	query := `
	SELECT
		to_json("t")
	FROM (
		SELECT
			"p"."id",
			"p"."title",
			"p"."description",
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
		WHERE "p"."id" = $1
		LIMIT 1
	) AS "t";`

	productBytes := make([]byte, 0)
	product := &products.Product{
		Images: make([]*entities.Images, 0),
	}

	// Query in bytes
	if err := r.db.Get(&productBytes, query, productId); err != nil {
		return nil, fmt.Errorf("get products failed: %v", err)
	}
	// Parse bytes to json object
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("unmarshal products failed: %v", err)
	}
	return product, nil
}
