package handlers

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/gofiber/fiber/v2"
)

type IMonitorHandler interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandler struct {
	Cfg config.IAppConfig
}

func MonitorHandler(cfg config.IAppConfig) IMonitorHandler {
	return &monitorHandler{
		Cfg: cfg,
	}
}

func (h *monitorHandler) HealthCheck(c *fiber.Ctx) error {
	return c.SendString("100%")
}
