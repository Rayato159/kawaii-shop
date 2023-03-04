package users_tests

import (
	"testing"

	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/users"
	"github.com/Rayato159/kawaii-shop/modules/users/repositories/patterns"
	"github.com/Rayato159/kawaii-shop/pkg/utils"
	kawaiitests "github.com/Rayato159/kawaii-shop/tests"
	"github.com/go-resty/resty/v2"
)

type testSignUpCustomer struct {
	req   *users.UserRegisterReq
	isErr bool
}

type testSignUpHandlerSuccess struct {
	url    string
	req    *users.UserRegisterReq
	expect int
}

type testSignUpHandlerError struct {
	url    string
	req    *users.UserRegisterReq
	expect int
}

func TestSignUpStruct(t *testing.T) {
	req := users.UserRegisterReq{
		Email:    "customer001kawaii.com",
		Username: "testcustomer001",
		Password: "123456",
	}
	if req.IsEmail() {
		t.Errorf("expect: %v, got: %v", false, req.IsEmail())
	}

	req.Email = "customer001@kawaii.com"
	if !req.IsEmail() {
		t.Errorf("expect: %v, got: %v", true, req.IsEmail())
	}
}

func TestHashing(t *testing.T) {
	password := "123456"
	req := users.UserRegisterReq{
		Email:    "customer001kawaii.com",
		Username: "testcustomer001",
		Password: password,
	}
	if err := req.BcryptHashing(); err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}
	if req.Password == password {
		t.Errorf("expect: %v, got: %v", false, req.Password == password)
	}
}

func TestSignUpCustomer(t *testing.T) {
	db := kawaiitests.Setup().GetDb()
	defer db.Close()

	tests := []testSignUpCustomer{
		{
			req: &users.UserRegisterReq{
				Email:    "customer001@kawaii.com",
				Username: "testcustomer001",
				Password: "123456",
			},
			isErr: true,
		},
		{
			req: &users.UserRegisterReq{
				Email:    "testcustomer001@kawaii.com",
				Username: "customer001",
				Password: "123456",
			},
			isErr: true,
		},
		{
			req: &users.UserRegisterReq{
				Email:    "cusotomer002@kawaii.com",
				Username: "customer002",
				Password: "123456",
			},
			isErr: false,
		},
	}

	for _, test := range tests {
		// Hashing password
		test.req.BcryptHashing()

		if test.isErr {
			_, err := patterns.InsertUser(db, test.req, false).Customer()
			if err == nil {
				t.Errorf("expect: %v, got: %v", "err", err)
			}
		} else {
			result, err := patterns.InsertUser(db, test.req, false).Customer()
			if err != nil {
				t.Errorf("expect: %v, got: %v", result, err)
			}

			user, err := result.Result()
			if err != nil {
				t.Errorf("expect: %v, got: %v", user, err)
			}
			if user == nil {
				t.Errorf("expect: %v, got: %v", "user", user)
			}
			utils.Debug(user)
		}
	}
}

func TestSignUpHandler(t *testing.T) {
	testsError := []testSignUpHandlerError{
		{
			url: "http://localhost:3000/v1/users/signup",
			req: &users.UserRegisterReq{
				Email:    "rainbowkawaii.com",
				Username: "rainbow",
				Password: "123456",
			},
			expect: 400,
		},
		{
			url: "http://localhost:3000/v1/users/signup",
			req: &users.UserRegisterReq{
				Email:    "rainbow@kawaii.com",
				Username: "customer001",
				Password: "123456",
			},
			expect: 400,
		},
		{
			url: "http://localhost:3000/v1/users/signup",
			req: &users.UserRegisterReq{
				Email:    "customer001@kawaii.com",
				Username: "rainbow",
				Password: "123456",
			},
			expect: 400,
		},
	}

	testsSuccess := []testSignUpHandlerSuccess{
		{
			url: "http://localhost:3000/v1/users/signup",
			req: &users.UserRegisterReq{
				Email:    "rainbow@kawaii.com",
				Username: "rainbow",
				Password: "123456",
			},
			expect: 201,
		},
	}

	// Init client
	client := resty.New()

	for _, req := range testsError {
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(req.req).
			SetResult(&entities.ErrorResponse{}).
			Post(req.url)
		if err != nil {
			t.Errorf("expect: %v, got: %v", nil, err)
		}
		if resp.StatusCode() != req.expect {
			t.Errorf("expect: %v, got: %v", req.expect, resp.StatusCode())
		}
	}

	for _, req := range testsSuccess {
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(req.req).
			SetResult(&users.UserPassport{}).
			Post(req.url)
		if err != nil {
			t.Errorf("expect: %v, got: %v", nil, err)
		}
		if resp.StatusCode() != req.expect {
			t.Errorf("expect: %v, got: %v", 201, resp.StatusCode())
		}
	}
}
