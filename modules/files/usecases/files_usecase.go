package usecases

import (
	filespkg "github.com/Rayato159/kawaii-shop/modules/files"
	"github.com/Rayato159/kawaii-shop/modules/files/repositories"
)

type IFilesUsecase interface {
	UploadFileTunnel(req *filespkg.FileReq) (*filespkg.FileRes, error)
}

type filesUsecase struct {
	filesRepository repositories.IFilesRepository
}

func FilesUsecase(fileRepository repositories.IFilesRepository) IFilesUsecase {
	return &filesUsecase{
		filesRepository: fileRepository,
	}
}

func (u *filesUsecase) UploadFileTunnel(req *filespkg.FileReq) (*filespkg.FileRes, error) {
	return nil, nil
}
