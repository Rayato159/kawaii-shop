package usecases

import "github.com/Rayato159/kawaii-shop/modules/middlewares/repositories"

type IMiddlewareUsecase interface{}

type middlewareUsecase struct {
	Repository repositories.IMiddlewareRepository
}

func MiddlewareUsecase(repo repositories.IMiddlewareRepository) IMiddlewareUsecase {
	return &middlewareUsecase{
		Repository: repo,
	}
}
