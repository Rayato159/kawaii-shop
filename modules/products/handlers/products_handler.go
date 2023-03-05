package handlers

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/products/usecases"
)

type IProductsHandler interface{}

type productsHandler struct {
	cfg             config.IConfig
	productsUsecase usecases.IProductsUsecase
}

func ProductHandler(cfg config.IConfig, productsUsecase usecases.IProductsUsecase) IProductsHandler {
	return &productsHandler{
		cfg:             cfg,
		productsUsecase: productsUsecase,
	}
}
