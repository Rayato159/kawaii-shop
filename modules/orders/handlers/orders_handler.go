package handlers

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/orders"
	_ordersUsecases "github.com/Rayato159/kawaii-shop/modules/orders/usecases"
	"github.com/gofiber/fiber/v2"
)

type ordersHandlerErrCode string

const (
	findOrdersErr ordersHandlerErrCode = "orders-001"
)

type IOrdersHandler interface {
	FindOrders(c *fiber.Ctx) error
}

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

func (h *ordersHandler) FindOrders(c *fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:     &entities.SortReq{},
		PaginateReq: &entities.PaginateReq{},
	}
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findOrdersErr),
			err.Error(),
		).Res()
	}
	// Paginate default
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	// Sort default
	if req.OrderBy == "" {
		req.OrderBy = "id"
	}
	if req.Sort == "" {
		req.Sort = "DESC"
	}

	orders := h.ordersUsecase.FindOrders(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, orders).Res()
}
