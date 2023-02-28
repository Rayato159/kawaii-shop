package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type server struct {
	app *fiber.App
	db  *sqlx.DB
	cfg config.IAppConfig
}

type IServer interface {
	Start()
	Db() *sqlx.DB
	Config() config.IAppConfig
}

func (s *server) Db() *sqlx.DB              { return s.db }
func (s *server) Config() config.IAppConfig { return s.cfg }

func (s *server) Start() {
	// Init Middleware
	middleware := InitMiddleware(s)
	middleware.Cors(s.app)

	// Init router
	v1 := s.app.Group("v1")

	// Import modules
	module := InitModule(v1, s, middleware)
	module.MonitorModule()
	module.OauthModule()

	// Log when server has started
	log.Printf("server is starting on %s", s.cfg.Url())

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("server is shutting down...")
		_ = s.app.Shutdown()
	}()

	// Listen to host:port
	s.app.Listen(s.cfg.Url())
}

func NewServer(cfg config.IAppConfig, db *sqlx.DB) IServer {
	return &server{
		app: fiber.New(fiber.Config{
			AppName:      cfg.Name(),
			BodyLimit:    cfg.BodyLimit(),
			ReadTimeout:  cfg.ReadTimeout(),
			WriteTimeout: cfg.WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
		db:  db,
		cfg: cfg,
	}
}
