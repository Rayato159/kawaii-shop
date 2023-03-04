package usecases

import (
	"github.com/Rayato159/kawaii-shop/modules/middlewares"
	"github.com/Rayato159/kawaii-shop/modules/middlewares/repositories"
)

type IMiddlewareUsecase interface {
	FindAccessToken(userId string, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewareUsecase struct {
	MiddlewareRepository repositories.IMiddlewareRepository
}

func MiddlewareUsecase(repo repositories.IMiddlewareRepository) IMiddlewareUsecase {
	return &middlewareUsecase{
		MiddlewareRepository: repo,
	}
}

func (u *middlewareUsecase) FindAccessToken(userId string, accessToken string) bool {
	return u.MiddlewareRepository.FindAccessToken(userId, accessToken)
}

func (u *middlewareUsecase) FindRole() ([]*middlewares.Role, error) {
	roles, err := u.MiddlewareRepository.FindRole()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
