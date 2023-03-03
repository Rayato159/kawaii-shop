package handlers

import (
	"strings"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/oauth"
	"github.com/Rayato159/kawaii-shop/modules/oauth/usecases"
	"github.com/gofiber/fiber/v2"
)

type oauthHandlerErrCode string

const (
	bodyParserErr     oauthHandlerErrCode = "oauth-001"
	signUpCustomerErr oauthHandlerErrCode = "oauth-002"
	getProfileErr     oauthHandlerErrCode = "oauth-003"
)

var oauthHandlerErrMsg = map[oauthHandlerErrCode]string{
	bodyParserErr:     "body parser failed",
	signUpCustomerErr: "insert customer error",
	getProfileErr:     "get profile error",
}

type IOauthHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
}

type oauthHandler struct {
	Cfg          config.IAppConfig
	OauthUsecase usecases.IOauthUsecase
}

func OauthHandler(cfg config.IAppConfig, usecase usecases.IOauthUsecase) IOauthHandler {
	return &oauthHandler{
		Cfg:          cfg,
		OauthUsecase: usecase,
	}
}

func (h *oauthHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request body parser
	req := new(oauth.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(bodyParserErr),
			oauthHandlerErrMsg[bodyParserErr],
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
	result, err := h.OauthUsecase.InsertCustomer(req)
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

func (h *oauthHandler) GetProfile(c *fiber.Ctx) error {
	// Set params
	userId := strings.Trim(c.Params("user_id"), " ")

	// Get profile
	result, err := h.OauthUsecase.GetProfile(userId)
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
