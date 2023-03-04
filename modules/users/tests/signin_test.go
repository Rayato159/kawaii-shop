package users_tests

import (
	"testing"

	"github.com/Rayato159/kawaii-shop/modules/users"
	"github.com/Rayato159/kawaii-shop/modules/users/repositories"
	"github.com/Rayato159/kawaii-shop/modules/users/usecases"
	"github.com/Rayato159/kawaii-shop/pkg/utils"
	kawaiitests "github.com/Rayato159/kawaii-shop/tests"
)

type testGetUserByEmail struct {
	email  string
	isErr  bool
	expect string
}

type testGetPassport struct {
	req    *users.UserCredential
	isErr  bool
	expect string
}

func TestGetUserCredentialCheck(t *testing.T) {
	db := kawaiitests.Setup().GetDb()
	tests := []testGetUserByEmail{
		{
			email:  "notfound@gmail.com",
			expect: "user not found",
			isErr:  true,
		},
		{
			email: "customer001@kawaii.com",
			expect: kawaiitests.ToJsonStringtify(&users.UserCredentialCheck{
				Id:       "U000001",
				Email:    "customer001@kawaii.com",
				Password: "$2a$10$8KzaNdKIMyOkASCH4QvSKuEMIY7Jc3vcHDuSJvXLii1rvBNgz60a6",
				Username: "customer001",
				Role:     "customer",
			}),
			isErr: false,
		},
	}

	usersRepo := repositories.UsersRepository(db)
	for _, req := range tests {
		if req.isErr {
			_, err := usersRepo.FindOneUserByEmail(req.email)
			if err == nil {
				t.Errorf("expect: %v, got: %v", "err", err)
				return
			}
			if err.Error() != req.expect {
				t.Errorf("expect: %v, got: %v", req.expect, err.Error())
			}
		} else {
			user, err := usersRepo.FindOneUserByEmail(req.email)
			if err != nil {
				t.Errorf("expect: %v, got: %v", nil, err)
			}
			if kawaiitests.ToJsonStringtify(user) != req.expect {
				t.Errorf("expect: %v, got: %v", req.expect, user)
			}
		}
	}
}

func TestGetPassport(t *testing.T) {
	init := kawaiitests.Setup()
	usersRepo := repositories.UsersRepository(init.GetDb())
	userUC := usecases.UsersUsecase(usersRepo, init.GetConfig())

	tests := []testGetPassport{
		{
			req: &users.UserCredential{
				Email:    "notfound101@kawaii.com",
				Password: "123456",
			},
			isErr:  true,
			expect: "user not found",
		},
		{
			req: &users.UserCredential{
				Email:    "customer001@kawaii.com",
				Password: "111111",
			},
			isErr:  true,
			expect: "password is invalid",
		},
		{
			req: &users.UserCredential{
				Email:    "customer001@kawaii.com",
				Password: "123456",
			},
			isErr: false,
		},
	}

	for _, req := range tests {
		if req.isErr {
			_, err := userUC.GetPassport(req.req)
			if err == nil {
				t.Errorf("expect: %v, got: %v", "err", err)
				return
			}
			if err.Error() != req.expect {
				t.Errorf("expect: %v, got: %v", req.expect, err.Error())
			}
		} else {
			res, err := userUC.GetPassport(req.req)
			if err != nil {
				t.Errorf("expect: %v, got: %v", nil, err)
			}
			utils.Debug(res)
		}
	}
}
