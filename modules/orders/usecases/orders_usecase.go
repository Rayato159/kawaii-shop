package usecases

import (
	_ordersRepositories "github.com/Rayato159/kawaii-shop/modules/orders/repositories"
)

type IOrdersUsecase interface{}

type ordersUsecase struct {
	ordersRepsotiory _ordersRepositories.IOrdersRepository
}

func OrdersUsecase(ordersRepsotiory _ordersRepositories.IOrdersRepository) IOrdersUsecase {
	return &ordersUsecase{
		ordersRepsotiory: ordersRepsotiory,
	}
}
