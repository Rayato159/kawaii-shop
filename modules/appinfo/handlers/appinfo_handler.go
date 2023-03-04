package handlers

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/appinfo/usecases"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/gofiber/fiber/v2"
)

type appinfoHandlerErrCode string

const (
	findCategoryErr appinfoHandlerErrCode = "app-001"
)

type IAppinfoHandler interface {
	FindCategory(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg            config.IConfig
	appinfoUsecase usecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, usecase usecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:            cfg,
		appinfoUsecase: usecase,
	}
}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	category, err := h.appinfoUsecase.FindCategory()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, category).Res()
}
