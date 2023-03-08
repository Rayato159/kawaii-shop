package product_tests

import (
	"fmt"
	"testing"

	"github.com/Rayato159/kawaii-shop/modules/appinfo"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/products"
	"github.com/Rayato159/kawaii-shop/modules/products/repositories"
	"github.com/Rayato159/kawaii-shop/pkg/utils"
	kawaiitests "github.com/Rayato159/kawaii-shop/tests"
)

type testFindProduct struct {
	req    *products.ProductFilter
	expect int
}

type testFindOneProduct struct {
	id    string
	isErr bool
}

type testAddProduct struct {
}

type testUpdateProduct struct {
	req *products.Product
}

func TestFindProduct(t *testing.T) {
	db := kawaiitests.Setup().GetDb()

	productsRepo := repositories.ProductsRepository(db)

	tests := []testFindProduct{
		{
			req: &products.ProductFilter{
				Id: "P000001",
				PaginateReq: &entities.PaginateReq{
					Page:  1,
					Limit: 10,
				},
				SortReq: &entities.SortReq{
					OrderBy: "product_id",
					Sort:    "DESC",
				},
			},
			expect: 1,
		},
		{
			req: &products.ProductFilter{
				Search: "fashion",
				PaginateReq: &entities.PaginateReq{
					Page:  1,
					Limit: 10,
				},
				SortReq: &entities.SortReq{
					OrderBy: "title",
					Sort:    "DESC",
				},
			},
			expect: 2,
		},
		{
			req: &products.ProductFilter{
				PaginateReq: &entities.PaginateReq{
					Page:  1,
					Limit: 10,
				},
				SortReq: &entities.SortReq{
					OrderBy: "product_id",
					Sort:    "DESC",
				},
			},
			expect: 6,
		},
		{
			req: &products.ProductFilter{
				Id: "P111111",
				PaginateReq: &entities.PaginateReq{
					Page:  1,
					Limit: 10,
				},
				SortReq: &entities.SortReq{
					OrderBy: "product_id",
					Sort:    "DESC",
				},
			},
			expect: 0,
		},
	}

	for _, test := range tests {
		products, count := productsRepo.FindProduct(test.req)
		if len(products) != test.expect {
			t.Errorf("expect: %v, got: %v", test.expect, len(products))
		}
		if count != test.expect {
			t.Errorf("expect: %v, got: %v", test.expect, count)
		}
		utils.Debug(products)
	}
}

func TestFindOneProduct(t *testing.T) {
	db := kawaiitests.Setup().GetDb()

	productsRepo := repositories.ProductsRepository(db)

	tests := []testFindOneProduct{
		{
			id:    "P111111",
			isErr: true,
		},
		{
			id:    "P000002",
			isErr: false,
		},
	}

	for _, test := range tests {
		if test.isErr {
			_, err := productsRepo.FindOneProduct(test.id)
			if err == nil {
				t.Errorf("expect: %v, got: %v", "err", err)
				return
			}
			fmt.Println(err)
		} else {
			product, err := productsRepo.FindOneProduct(test.id)
			if err != nil {
				t.Errorf("expect: %v, got: %v", nil, err)
			}
			utils.Debug(product)
		}
	}
}

func TestAddProduct(t *testing.T) {
	db := kawaiitests.Setup().GetDb()
	productsRepo := repositories.ProductsRepository(db)

	product, err := productsRepo.InsertProduct(&products.Product{
		Title:       "air purifier",
		Description: "something I don't know",
		Category: &appinfo.Category{
			Id: 3,
		},
		Images: []*entities.Images{
			{
				FileName: utils.RandomFileName("jpg"),
				Url:      "https://i.pinimg.com/564x/b1/8e/1b/b18e1b87f50ef7ccd1e2c58c408fe35e.jpg",
			},
			{
				FileName: utils.RandomFileName("jpg"),
				Url:      "https://i.pinimg.com/564x/b1/8e/1b/b18e1b87f50ef7ccd1e2c58c408fe35e.jpg",
			},
			{
				FileName: utils.RandomFileName("jpg"),
				Url:      "https://i.pinimg.com/564x/b1/8e/1b/b18e1b87f50ef7ccd1e2c58c408fe35e.jpg",
			},
		},
	})
	if err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}
	utils.Debug(product)
}

func TestUpdateProductRepo(t *testing.T) {
	db := kawaiitests.Setup().GetDb()
	productsRepo := repositories.ProductsRepository(db)

	tests := []testUpdateProduct{
		{
			req: &products.Product{
				Id:          "P000006",
				Title:       "Disc",
				Description: "Just a music disc",
				Category:    &appinfo.Category{},
				Images:      make([]*entities.Images, 0),
			},
		},
		{
			req: &products.Product{
				Id:          "P000006",
				Title:       "",
				Description: "Hello World!",
				Category:    &appinfo.Category{},
				Images:      make([]*entities.Images, 0),
			},
		},
	}

	for _, test := range tests {
		product, err := productsRepo.UpdateProduct(test.req)
		if err != nil {
			t.Errorf("expect: %v, got: %v", nil, err)
		}
		utils.Debug(product)
	}
}
