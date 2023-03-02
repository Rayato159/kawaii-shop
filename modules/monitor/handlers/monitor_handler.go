package handlers

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/monitor"
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
	res := &monitor.Monitor{
		Name:    h.Cfg.Name(),
		Version: h.Cfg.Version(),
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}
