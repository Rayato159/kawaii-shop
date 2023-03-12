package usecases

import (
	"math"

	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/orders"
	_ordersRepositories "github.com/Rayato159/kawaii-shop/modules/orders/repositories"
)

type IOrdersUsecase interface {
	FindOrder(req *orders.OrderFilter) *entities.PaginateRes
	FindOneOrder(orderId string) (*orders.Order, error)
}

type ordersUsecase struct {
	ordersRepsotiory _ordersRepositories.IOrdersRepository
}

func OrdersUsecase(ordersRepsotiory _ordersRepositories.IOrdersRepository) IOrdersUsecase {
	return &ordersUsecase{
		ordersRepsotiory: ordersRepsotiory,
	}
}

func (u *ordersUsecase) FindOrder(req *orders.OrderFilter) *entities.PaginateRes {
	orders, count := u.ordersRepsotiory.FindOrder(req)

	return &entities.PaginateRes{
		Data:      orders,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *ordersUsecase) FindOneOrder(orderId string) (*orders.Order, error) {
	order, err := u.ordersRepsotiory.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}
