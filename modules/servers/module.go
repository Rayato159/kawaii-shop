package servers

import (
	_middlewareHandlers "github.com/Rayato159/kawaii-shop/modules/middlewares/handlers"
	_middlewareRepositories "github.com/Rayato159/kawaii-shop/modules/middlewares/repositories"
	_middlewareUsecases "github.com/Rayato159/kawaii-shop/modules/middlewares/usecases"

	_monitorHandlers "github.com/Rayato159/kawaii-shop/modules/monitor/handlers"

	_usersHandlers "github.com/Rayato159/kawaii-shop/modules/users/handlers"
	_usersRepositories "github.com/Rayato159/kawaii-shop/modules/users/repositories"
	_usersUsecases "github.com/Rayato159/kawaii-shop/modules/users/usecases"

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
	UsersModule()
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
	f.router.Get("/", _monitorHandlers.MonitorHandler(f.server.cfg.App()).HealthCheck)
}

func (f *ModuleFactory) UsersModule() {
	repository := _usersRepositories.UsersRepository(f.server.Db())
	usecase := _usersUsecases.UsersUsecase(repository, f.server.cfg)
	handler := _usersHandlers.UsersHandler(f.server.cfg.App(), usecase)

	router := f.router.Group("/users")
	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/signout", handler.SignOut)
	router.Get("/:user_id", handler.GetProfile)
}
