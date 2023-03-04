package usecases

import (
	"github.com/Rayato159/kawaii-shop/modules/appinfo"
	"github.com/Rayato159/kawaii-shop/modules/appinfo/repositories"
)

type IAppinfoUsecase interface {
	FindCategory() ([]*appinfo.Category, error)
}

type appinfoUsecase struct {
	appinfoRepository repositories.IAppinfoRepository
}

func AppinfoUsecase(repo repositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: repo,
	}
}

func (u *appinfoUsecase) FindCategory() ([]*appinfo.Category, error) {
	category, err := u.appinfoRepository.FindCategory()
	if err != nil {
		return nil, err
	}
	return category, nil
}
