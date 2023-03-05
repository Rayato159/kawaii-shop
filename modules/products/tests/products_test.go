package product_tests

import (
	"fmt"
	"testing"

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

func TestFindProduct(t *testing.T) {
	db := kawaiitests.Setup().GetDb()

	productsRepo := repositories.ProductRepository(db)

	tests := []testFindProduct{
		{
			req: &products.ProductFilter{
				Id: "P000001",
			},
			expect: 1,
		},
		{
			req: &products.ProductFilter{
				Search: "fashion",
			},
			expect: 2,
		},
		{
			req:    &products.ProductFilter{},
			expect: 6,
		},
		{
			req: &products.ProductFilter{
				Id: "P111111",
			},
			expect: 0,
		},
	}

	for _, test := range tests {
		products := productsRepo.FindProduct(test.req)
		if len(products) != test.expect {
			t.Errorf("expect: %v, got: %v", test.expect, len(products))
		}
		utils.Debug(products)
	}
}

func TestFindOneProduct(t *testing.T) {
	db := kawaiitests.Setup().GetDb()

	productsRepo := repositories.ProductRepository(db)

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
