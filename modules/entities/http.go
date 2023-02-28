package entities

import (
	"github.com/Rayato159/kawaii-shop/pkg/kawaiilogger"
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Success(code int, data any) IResponse
	Error(code int, traceId, msg string) IResponse
	Res() error
}

func (r *Response) Success(code int, data any) IResponse {
	r.StatusCode = code
	r.Data = data
	kawaiilogger.InitKawaiiLogger(r.Context, &r.Data).Print().Save()
	return r
}

func (r *Response) Error(code int, traceId, msg string) IResponse {
	r.StatusCode = code
	r.ErrorRes = &ErrorResponse{
		TraceId: traceId,
		Msg:     msg,
	}
	r.IsError = true
	kawaiilogger.InitKawaiiLogger(r.Context, &r.ErrorRes).Print().Save()
	return r
}

func (r *Response) Res() error {
	if r.IsError {
		return r.Context.Status(r.StatusCode).JSON(&r.ErrorRes)
	}
	return r.Context.Status(r.StatusCode).JSON(&r.Data)
}

type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	Context    *fiber.Ctx
	IsError    bool
}

type ErrorResponse struct {
	TraceId string `json:"trace_id"`
	Msg     string `json:"message"`
}

func NewResponse(c *fiber.Ctx) IResponse {
	return &Response{
		Context: c,
	}
}
