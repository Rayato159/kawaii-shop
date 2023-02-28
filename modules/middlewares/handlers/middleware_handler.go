package handlers

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/middlewares/usecases"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type IMiddlewareHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	JwtAuth() fiber.Handler
}

type middlewareHandler struct {
	Cfg     config.IAppConfig
	Usecase usecases.IMiddlewareUsecase
}

func MiddlewareHandler(cfg config.IAppConfig, usecase usecases.IMiddlewareUsecase) IMiddlewareHandler {
	return &middlewareHandler{
		Cfg:     cfg,
		Usecase: usecase,
	}
}

func (h *middlewareHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewareHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

func (h *middlewareHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			"r-001",
			"router not found",
		).Res()
	}
}
