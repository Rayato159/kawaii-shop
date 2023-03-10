package usecases

import (
	"math"

	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/products"
	"github.com/Rayato159/kawaii-shop/modules/products/repositories"
)

type IProductsUsecase interface {
	FindProduct(req *products.ProductFilter) *entities.PaginateRes
	FindOneProduct(productId string) (*products.Product, error)
	AddProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
	UpdateProduct(req *products.Product) (*products.Product, error)
}

type productsUsecase struct {
	productRepository repositories.IProductsRepository
}

func ProductsUsecase(productsRepo repositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productRepository: productsRepo,
	}
}

func (u *productsUsecase) FindProduct(req *products.ProductFilter) *entities.PaginateRes {
	products, count := u.productRepository.FindProduct(req)

	return &entities.PaginateRes{
		Data:      products,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *productsUsecase) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.productRepository.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productsUsecase) AddProduct(req *products.Product) (*products.Product, error) {
	product, err := u.productRepository.InsertProduct(req)
	if err != nil {
		return nil, err
	}
	return product, err
}

func (u *productsUsecase) UpdateProduct(req *products.Product) (*products.Product, error) {
	product, err := u.productRepository.UpdateProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productsUsecase) DeleteProduct(productId string) error {
	if err := u.productRepository.DeleteProduct(productId); err != nil {
		return err
	}
	return nil
}
