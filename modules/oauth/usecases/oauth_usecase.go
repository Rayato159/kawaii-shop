package usecases

import "github.com/Rayato159/kawaii-shop/modules/oauth/repositories"

type IOauthUsecase interface{}

type oauthUsecase struct {
	Repository repositories.IOauthRepository
}

func OauthUsecase(repo repositories.IOauthRepository) IOauthUsecase {
	return &oauthUsecase{
		Repository: repo,
	}
}
