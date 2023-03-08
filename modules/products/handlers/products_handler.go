package handlers

import (
	"strings"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/appinfo"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/products"
	"github.com/Rayato159/kawaii-shop/modules/products/usecases"
	"github.com/gofiber/fiber/v2"
)

type productsHandlerErrCode string

const (
	findProductErr    productsHandlerErrCode = "products-001"
	findOneProductErr productsHandlerErrCode = "products-002"
	addProductErr     productsHandlerErrCode = "products-003"
)

type IProductsHandler interface {
	FindProduct(c *fiber.Ctx) error
	FindOneProduct(c *fiber.Ctx) error
	AddProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type productsHandler struct {
	cfg             config.IConfig
	productsUsecase usecases.IProductsUsecase
}

func ProductsHandler(cfg config.IConfig, productsUsecase usecases.IProductsUsecase) IProductsHandler {
	return &productsHandler{
		cfg:             cfg,
		productsUsecase: productsUsecase,
	}
}

func (h *productsHandler) FindProduct(c *fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginateReq: &entities.PaginateReq{},
		SortReq:     &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findProductErr),
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
		req.OrderBy = "title"
	}
	if req.Sort == "" {
		req.Sort = "ASC"
	}

	products := h.productsUsecase.FindProduct(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, products).Res()
}

func (h *productsHandler) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productsHandler) AddProduct(c *fiber.Ctx) error {
	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Images, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addProductErr),
			err.Error(),
		).Res()
	}
	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addProductErr),
			"category id is invalid",
		).Res()
	}

	product, err := h.productsUsecase.AddProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(addProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, product).Res()
}

func (h *productsHandler) DeleteProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	if err := h.productsUsecase.DeleteProduct(productId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(addProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
