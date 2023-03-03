package usecases

import (
	"github.com/Rayato159/kawaii-shop/modules/oauth"
	"github.com/Rayato159/kawaii-shop/modules/oauth/repositories"
)

type IOauthUsecase interface {
	InsertCustomer(req *oauth.UserRegisterReq) (*oauth.UserPassport, error)
	GetProfile(userId string) (*oauth.User, error)
}

type oauthUsecase struct {
	OauthRepository repositories.IOauthRepository
}

func OauthUsecase(repo repositories.IOauthRepository) IOauthUsecase {
	return &oauthUsecase{
		OauthRepository: repo,
	}
}

func (u *oauthUsecase) InsertCustomer(req *oauth.UserRegisterReq) (*oauth.UserPassport, error) {
	// Hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Inserting user
	result, err := u.OauthRepository.InsertCustomer(req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *oauthUsecase) GetProfile(userId string) (*oauth.User, error) {
	profile, err := u.OauthRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
