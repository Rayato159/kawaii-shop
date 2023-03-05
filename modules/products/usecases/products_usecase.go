package usecases

import "github.com/Rayato159/kawaii-shop/modules/products/repositories"

type IProductsUsecase interface{}

type productsUsecase struct {
	productRepository repositories.IProductsRepository
}

func ProductUsecase(productsRepo repositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productRepository: productsRepo,
	}
}
