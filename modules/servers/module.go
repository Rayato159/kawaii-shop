package servers

import (
	_middlewareHandlers "github.com/Rayato159/kawaii-shop/modules/middlewares/handlers"
	_middlewareRepositories "github.com/Rayato159/kawaii-shop/modules/middlewares/repositories"
	_middlewareUsecases "github.com/Rayato159/kawaii-shop/modules/middlewares/usecases"

	_monitorHandlers "github.com/Rayato159/kawaii-shop/modules/monitor/handlers"

	_oauthHandlers "github.com/Rayato159/kawaii-shop/modules/oauth/handlers"
	_oauthRepositories "github.com/Rayato159/kawaii-shop/modules/oauth/repositories"
	_oauthUsecases "github.com/Rayato159/kawaii-shop/modules/oauth/usecases"

	"github.com/gofiber/fiber/v2"
)

// Middleware
func InitMiddleware(s *server) _middlewareHandlers.IMiddlewareHandler {
	repository := _middlewareRepositories.MiddlewareRepository(s.db)
	usecase := _middlewareUsecases.MiddlewareUsecase(repository)
	handler := _middlewareHandlers.MiddlewareHandler(s.cfg, usecase)
	return handler
}

// Module
type IModuleFactory interface {
	MonitorModule()
	OauthModule()
}

type ModuleFactory struct {
	router     fiber.Router
	server     *server
	middleware _middlewareHandlers.IMiddlewareHandler
}

func InitModule(r fiber.Router, s *server, m _middlewareHandlers.IMiddlewareHandler) IModuleFactory {
	return &ModuleFactory{
		router:     r,
		server:     s,
		middleware: m,
	}
}

func (f *ModuleFactory) MonitorModule() {
	f.router.Get("/", _monitorHandlers.MonitorHandler(f.server.Config()).HealthCheck)
}

func (f *ModuleFactory) OauthModule() {
	repository := _oauthRepositories.OauthRepository(f.server.Db())
	usecase := _oauthUsecases.OauthUsecase(repository)
	handler := _oauthHandlers.OauthHandler(f.server.Config(), usecase)
	_ = handler
}
