package handlers

import (
	"strings"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/users"
	"github.com/Rayato159/kawaii-shop/modules/users/usecases"
	"github.com/Rayato159/kawaii-shop/pkg/kawaiiauth"
	"github.com/gofiber/fiber/v2"
)

type usersHandlerErrCode string

const (
	bodyParserErr         usersHandlerErrCode = "users-001"
	signUpCustomerErr     usersHandlerErrCode = "users-002"
	getProfileErr         usersHandlerErrCode = "users-003"
	signInErr             usersHandlerErrCode = "users-004"
	refreshPassportErr    usersHandlerErrCode = "users-005"
	signOutErr            usersHandlerErrCode = "users-006"
	generateAdminTokenErr usersHandlerErrCode = "users-007"
)

var usersHandlerErrMsg = map[usersHandlerErrCode]string{
	bodyParserErr:         "body parser failed",
	signUpCustomerErr:     "insert customer error",
	getProfileErr:         "get profile error",
	signInErr:             "sign in error",
	refreshPassportErr:    "refresh token error",
	signOutErr:            "sign out error",
	generateAdminTokenErr: "generate admin token error",
}

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
}

type usersHandler struct {
	Cfg          config.IConfig
	UsersUsecase usecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, usecase usecases.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		Cfg:          cfg,
		UsersUsecase: usecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(bodyParserErr),
			usersHandlerErrMsg[bodyParserErr],
		).Res()
	}

	// Email validatio
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(bodyParserErr),
			"email pattern is invalid",
		).Res()
	}

	// Insert
	result, err := h.UsersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username have been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email have been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(bodyParserErr),
			usersHandlerErrMsg[bodyParserErr],
		).Res()
	}

	passport, err := h.UsersUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(bodyParserErr),
			usersHandlerErrMsg[bodyParserErr],
		).Res()
	}

	passport, err := h.UsersUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(bodyParserErr),
			usersHandlerErrMsg[bodyParserErr],
		).Res()
	}

	if err := h.UsersUsecase.DeleteOauth(req.Code); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) GetProfile(c *fiber.Ctx) error {
	// Set params
	userId := strings.Trim(c.Params("user_id"), " ")

	// Get profile
	result, err := h.UsersUsecase.GetProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user profile failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(getProfileErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(getProfileErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *usersHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := kawaiiauth.NewKawaiiAuth(
		kawaiiauth.Admin,
		h.Cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateAdminTokenErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, &struct {
		Token string `json:"token"`
	}{
		Token: adminToken.SignToken(),
	},
	).Res()
}
