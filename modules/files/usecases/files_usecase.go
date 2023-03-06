package usecases

import "github.com/Rayato159/kawaii-shop/modules/files/repositories"

type IFilesUsecase interface{}

type filesUsecase struct {
	filesRepository repositories.IFilesRepository
}

func FilesUsecase(fileRepository repositories.IFilesRepository) IFilesUsecase {
	return &filesUsecase{
		filesRepository: fileRepository,
	}
}
