package modules

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/monitor/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func MonitorModule(r fiber.Router, cfg config.IAppConfig) {
	r.Get("/", handlers.MonitorHandler(cfg).HealthCheck)
}

func OauthModule(r fiber.Router, cfg config.IAppConfig, db *sqlx.DB) {
}
