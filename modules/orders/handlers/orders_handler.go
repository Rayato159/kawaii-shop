package handlers

import (
	"strings"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/orders"
	_ordersUsecases "github.com/Rayato159/kawaii-shop/modules/orders/usecases"
	"github.com/gofiber/fiber/v2"
)

type ordersHandlerErrCode string

const (
	findOrderErr    ordersHandlerErrCode = "orders-001"
	findOneOrderErr ordersHandlerErrCode = "orders-002"
)

type IOrdersHandler interface {
	FindOrder(c *fiber.Ctx) error
	FindOneOrder(c *fiber.Ctx) error
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

func (h *ordersHandler) FindOrder(c *fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:     &entities.SortReq{},
		PaginateReq: &entities.PaginateReq{},
	}
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findOrderErr),
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

	orders := h.ordersUsecase.FindOrder(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, orders).Res()
}

func (h *ordersHandler) FindOneOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")

	order, err := h.ordersUsecase.FindOneOrder(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findOneOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}
