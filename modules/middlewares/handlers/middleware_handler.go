package handlers

import (
	"fmt"
	"strings"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/middlewares/usecases"
	"github.com/Rayato159/kawaii-shop/pkg/kawaiiauth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandlerErrCode string

const (
	jwtAuthErr     middlewareHandlerErrCode = "middleware-001"
	paramsCheckErr middlewareHandlerErrCode = "middleware-002"
)

type IMiddlewareHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
}

type middlewareHandler struct {
	Cfg               config.IConfig
	MiddlewareUsecase usecases.IMiddlewareUsecase
}

func MiddlewareHandler(cfg config.IConfig, usecase usecases.IMiddlewareUsecase) IMiddlewareHandler {
	return &middlewareHandler{
		Cfg:               cfg,
		MiddlewareUsecase: usecase,
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
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := kawaiiauth.ParseToken(h.Cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		fmt.Println(h.MiddlewareUsecase.FindAccessToken(claims.Id, token))
		if !h.MiddlewareUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}

		// Set userId
		c.Locals("userId", claims.Id)
		return c.Next()
	}
}

func (h *middlewareHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId").(string)
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(jwtAuthErr),
				"never gonna give you up",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewareHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "01/02/2006",
		TimeZone:   "Bangkok/Asia",
	})
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
