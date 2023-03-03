package oauth_tests

import (
	"testing"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/oauth"
	"github.com/Rayato159/kawaii-shop/modules/oauth/repositories/patterns"
	"github.com/Rayato159/kawaii-shop/pkg/databases"
	"github.com/Rayato159/kawaii-shop/pkg/utils"
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
)

type testSignUpCustomer struct {
	req   *oauth.UserRegisterReq
	isErr bool
}

type testSignUpHandlerSuccess struct {
	url    string
	req    *oauth.UserRegisterReq
	expect int
}

type testSignUpHandlerError struct {
	url    string
	req    *oauth.UserRegisterReq
	expect int
}

type testOauthConfig struct {
	cfg config.IConfig
}

type ITestOauthConfig interface {
	getDb() *sqlx.DB
}

func setup() ITestOauthConfig {
	cfg := config.LoadConfig("../../../.env.test")
	return &testOauthConfig{
		cfg: cfg,
	}
}

func (cfg *testOauthConfig) getDb() *sqlx.DB {
	return databases.DbConnect(cfg.cfg.Db())
}

func TestSignUpStruct(t *testing.T) {
	req := oauth.UserRegisterReq{
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
	req := oauth.UserRegisterReq{
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
	db := setup().getDb()
	defer db.Close()

	tests := []testSignUpCustomer{
		{
			req: &oauth.UserRegisterReq{
				Email:    "customer001@kawaii.com",
				Username: "testcustomer001",
				Password: "123456",
			},
			isErr: true,
		},
		{
			req: &oauth.UserRegisterReq{
				Email:    "testcustomer001@kawaii.com",
				Username: "customer001",
				Password: "123456",
			},
			isErr: true,
		},
		{
			req: &oauth.UserRegisterReq{
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
			_, err := patterns.InsertUser(db, test.req).Customer()
			if err == nil {
				t.Errorf("expect: %v, got: %v", "err", err)
			}
		} else {
			result, err := patterns.InsertUser(db, test.req).Customer()
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
			url: "http://localhost:3000/v1/oauth/signup",
			req: &oauth.UserRegisterReq{
				Email:    "rainbowkawaii.com",
				Username: "rainbow",
				Password: "123456",
			},
			expect: 400,
		},
		{
			url: "http://localhost:3000/v1/oauth/signup",
			req: &oauth.UserRegisterReq{
				Email:    "rainbow@kawaii.com",
				Username: "customer001",
				Password: "123456",
			},
			expect: 400,
		},
		{
			url: "http://localhost:3000/v1/oauth/signup",
			req: &oauth.UserRegisterReq{
				Email:    "customer001@kawaii.com",
				Username: "rainbow",
				Password: "123456",
			},
			expect: 400,
		},
	}

	testsSuccess := []testSignUpHandlerSuccess{
		{
			url: "http://localhost:3000/v1/oauth/signup",
			req: &oauth.UserRegisterReq{
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
			SetResult(&oauth.UserPassport{}).
			Post(req.url)
		if err != nil {
			t.Errorf("expect: %v, got: %v", nil, err)
		}
		if resp.StatusCode() != req.expect {
			t.Errorf("expect: %v, got: %v", 201, resp.StatusCode())
		}
	}
}
