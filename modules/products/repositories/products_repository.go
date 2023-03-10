package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/products"
	"github.com/Rayato159/kawaii-shop/modules/products/repositories/patterns"

	_filesUsecases "github.com/Rayato159/kawaii-shop/modules/files/usecases"

	"github.com/jmoiron/sqlx"
)

type IProductsRepository interface {
	FindProduct(req *products.ProductFilter) ([]*products.Product, int)
	FindOneProduct(productId string) (*products.Product, error)
	InsertProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
	UpdateProduct(req *products.Product) (*products.Product, error)
}

type productsRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase _filesUsecases.IFilesUsecase
}

func ProductsRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase _filesUsecases.IFilesUsecase) IProductsRepository {
	return &productsRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *productsRepository) FindProduct(req *products.ProductFilter) ([]*products.Product, int) {
	builder := patterns.FindProductBuilder(r.db, req)
	engineer := patterns.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()
	return result, count
}

func (r *productsRepository) FindOneProduct(productId string) (*products.Product, error) {
	query := `
	SELECT
		to_json("t")
	FROM (
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

func (r *productsRepository) InsertProduct(req *products.Product) (*products.Product, error) {
	builder := patterns.InsertProductBuilder(r.db, req)
	productId, err := patterns.InsertProductEngineer(builder).InsertProduct()
	if err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productsRepository) UpdateProduct(req *products.Product) (*products.Product, error) {
	builder := patterns.UpdateProductBuilder(r.db, req, r.filesUsecase)
	engineer := patterns.UpdateProductEngineer(builder)

	if err := engineer.UpdateProduct(); err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(req.Id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productsRepository) DeleteProduct(productId string) error {
	query := `
	DELETE FROM "products" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, productId); err != nil {
		return fmt.Errorf("delete products failed: %v", err)
	}
	return nil
}
