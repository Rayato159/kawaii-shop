package modules

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/monitor/handlers"
	"github.com/gofiber/fiber/v2"
)

func Monitor(r fiber.Router, cfg config.IAppConfig) {
	r.Get("/", handlers.MonitorHandler(cfg).HealthCheck)
}
