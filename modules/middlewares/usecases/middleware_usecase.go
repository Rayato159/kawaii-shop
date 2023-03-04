package usecases

import "github.com/Rayato159/kawaii-shop/modules/middlewares/repositories"

type IMiddlewareUsecase interface {
	FindAccessToken(userId string, accessToken string) bool
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
