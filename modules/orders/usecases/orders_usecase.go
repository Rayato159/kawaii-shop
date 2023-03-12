package usecases

import (
	"fmt"
	"math"

	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/orders"
	_ordersRepositories "github.com/Rayato159/kawaii-shop/modules/orders/repositories"
	_productsRepositories "github.com/Rayato159/kawaii-shop/modules/products/repositories"
)

type IOrdersUsecase interface {
	FindOrder(req *orders.OrderFilter) *entities.PaginateRes
	FindOneOrder(orderId string) (*orders.Order, error)
	InsertOrder(req *orders.Order) (*orders.Order, error)
	UpdateOrder(req *orders.UpdateOrderReq) (*orders.Order, error)
}

type ordersUsecase struct {
	ordersRepsotiory   _ordersRepositories.IOrdersRepository
	productsRepsotiory _productsRepositories.IProductsRepository
}

func OrdersUsecase(ordersRepsotiory _ordersRepositories.IOrdersRepository, productsRepsotiory _productsRepositories.IProductsRepository) IOrdersUsecase {
	return &ordersUsecase{
		ordersRepsotiory:   ordersRepsotiory,
		productsRepsotiory: productsRepsotiory,
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

func (u *ordersUsecase) InsertOrder(req *orders.Order) (*orders.Order, error) {
	// Search product if exists
	for i := range req.Products {
		if req.Products[i].Product == nil {
			return nil, fmt.Errorf("product is nil")
		}
		prod, err := u.productsRepsotiory.FindOneProduct(req.Products[i].Product.Id)
		if err != nil {
			return nil, err
		}
		req.TotalPaid += req.Products[i].Product.Price * float64(req.Products[i].Qty)
		req.Products[i].Product = prod
	}

	orderId, err := u.ordersRepsotiory.InsertOrder(req)
	if err != nil {
		return nil, err
	}

	order, err := u.ordersRepsotiory.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (u *ordersUsecase) UpdateOrder(req *orders.UpdateOrderReq) (*orders.Order, error) {
	if err := u.ordersRepsotiory.UpdateOrder(req); err != nil {
		return nil, err
	}

	order, err := u.ordersRepsotiory.FindOneOrder(req.OrderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}
