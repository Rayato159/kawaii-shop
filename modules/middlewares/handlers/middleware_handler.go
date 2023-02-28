package handlers

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/middlewares/usecases"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type IMiddlewareHandler interface {
	JwtAuth() fiber.Handler
	Cors(a *fiber.App) fiber.Handler
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

func (h *middlewareHandler) Cors(a *fiber.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		a.Use(cors.New(cors.Config{
			Next:             nil,
			AllowOrigins:     "*",
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
			AllowHeaders:     "",
			AllowCredentials: false,
			ExposeHeaders:    "",
			MaxAge:           0,
		}))
		return c.Next()
	}
}

func (h *middlewareHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
