package handlers

import (
	"github.com/Rayato159/kawaii-shop/config"
	_ordersUsecases "github.com/Rayato159/kawaii-shop/modules/orders/usecases"
)

type IOrdersHandler interface{}

type ordersHandler struct {
	cfg           config.IConfig
	ordersUsecase _ordersUsecases.IOrdersUsecase
}

func OrdersHandler(cfg config.IConfig, ordersUsecase _ordersUsecases.IOrdersUsecase) IOrdersHandler {
	return &ordersHandler{
		cfg:           cfg,
		ordersUsecase: ordersUsecase,
	}
}
