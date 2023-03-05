package usecases

import (
	"github.com/Rayato159/kawaii-shop/modules/appinfo"
	"github.com/Rayato159/kawaii-shop/modules/appinfo/repositories"
)

type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) ([]*appinfo.Category, error)
	DeleteCategory(categoryId int) error
}

type appinfoUsecase struct {
	appinfoRepository repositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepo repositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepo,
	}
}

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	category, err := u.appinfoRepository.FindCategory(req)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) ([]*appinfo.Category, error) {
	if err := u.appinfoRepository.InsertCategory(req); err != nil {
		return nil, err
	}
	category, err := u.appinfoRepository.FindCategory(&appinfo.CategoryFilter{})
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (u *appinfoUsecase) DeleteCategory(categoryId int) error {
	if err := u.appinfoRepository.DeleteCategory(categoryId); err != nil {
		return err
	}
	return nil
}
