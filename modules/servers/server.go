package servers

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules"
	"github.com/gofiber/fiber/v2"
)

type server struct {
	app *fiber.App
	cfg config.IAppConfig
}

type IServer interface {
	Start()
}

func (s *server) Start() {
	v1 := s.app.Group("v1")

	modules.Monitor(v1, s.cfg)

	s.app.Listen(s.cfg.Url())
}

func NewServer(cfg config.IAppConfig) IServer {
	return &server{
		app: fiber.New(fiber.Config{}),
		cfg: cfg,
	}
}
