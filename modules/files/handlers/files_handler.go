package handlers

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/files/usecases"
)

type IFilesHandler interface{}

type filesHandler struct {
	cfg          config.IConfig
	filesUsecase usecases.IFilesUsecase
}

func FilesHandler(cfg config.IConfig, filesUsecase usecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}
