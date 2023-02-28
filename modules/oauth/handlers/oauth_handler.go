package handlers

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/oauth/usecases"
)

type IOauthHandler interface {
}

type oauthHandler struct {
	Cfg     config.IAppConfig
	Usecase usecases.IOauthUsecase
}

func OauthModule(cfg config.IAppConfig, usecase usecases.IOauthUsecase) IOauthHandler {
	return &oauthHandler{
		Cfg:     cfg,
		Usecase: usecase,
	}
}
