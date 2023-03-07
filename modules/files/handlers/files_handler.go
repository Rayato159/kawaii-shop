package handlers

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	filespkg "github.com/Rayato159/kawaii-shop/modules/files"
	"github.com/Rayato159/kawaii-shop/modules/files/usecases"
	"github.com/Rayato159/kawaii-shop/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type filesHandlerErrCode string

const (
	uploadFileErr filesHandlerErrCode = "files-001"
)

type IFilesHandler interface {
	UploadFiles(c *fiber.Ctx) error
}

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

func (h *filesHandler) UploadFiles(c *fiber.Ctx) error {
	// Init req obj
	req := make([]*filespkg.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadFileErr),
			err.Error(),
		).Res()
	}
	files := form.File["files"]
	destination := c.FormValue("destination")

	// Files validation
	extensionMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}
	for _, file := range files {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extensionMap[ext] != ext {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadFileErr),
				"extension is not acceptable",
			).Res()
		}

		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadFileErr),
				fmt.Sprintf("file size must less than %d MiB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		req = append(req, &filespkg.FileReq{
			File:        file,
			Destination: destination,
			FileName:    utils.RandomFileName(ext),
			Extension:   ext,
		})
	}
	utils.Debug(req)

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
