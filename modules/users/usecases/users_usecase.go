package usecases

import (
	"github.com/Rayato159/kawaii-shop/modules/users"
	"github.com/Rayato159/kawaii-shop/modules/users/repositories"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetProfile(userId string) (*users.User, error)
}

type usersUsecase struct {
	UsersRepository repositories.IUsersRepository
}

func UsersUsecase(repo repositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
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
