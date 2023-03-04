package usecases

import (
	"fmt"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/users"
	"github.com/Rayato159/kawaii-shop/modules/users/repositories"
	"github.com/Rayato159/kawaii-shop/pkg/kawaiiauth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetProfile(userId string) (*users.User, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	UsersRepository repositories.IUsersRepository
}

func UsersUsecase(repo repositories.IUsersRepository, cfg config.IConfig) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		UsersRepository: repo,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// Hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Inserting user
	result, err := u.UsersRepository.InsertCustomer(req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *usersUsecase) GetProfile(userId string) (*users.User, error) {
	profile, err := u.UsersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	// Find user
	user, err := u.UsersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	// Generate token
	accessToken, err := kawaiiauth.NewKawaiiAuth(kawaiiauth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := kawaiiauth.NewKawaiiAuth(kawaiiauth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
	})
	if err != nil {
		return nil, err
	}

	// Set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	// Insert oauth table
	if err := u.UsersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil
}
