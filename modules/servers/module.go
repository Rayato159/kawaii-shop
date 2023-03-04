package servers

import (
	_middlewareHandlers "github.com/Rayato159/kawaii-shop/modules/middlewares/handlers"
	_middlewareRepositories "github.com/Rayato159/kawaii-shop/modules/middlewares/repositories"
	_middlewareUsecases "github.com/Rayato159/kawaii-shop/modules/middlewares/usecases"

	_monitorHandlers "github.com/Rayato159/kawaii-shop/modules/monitor/handlers"

	_usersHandlers "github.com/Rayato159/kawaii-shop/modules/users/handlers"
	_usersRepositories "github.com/Rayato159/kawaii-shop/modules/users/repositories"
	_usersUsecases "github.com/Rayato159/kawaii-shop/modules/users/usecases"

	_appinfoHandlers "github.com/Rayato159/kawaii-shop/modules/appinfo/handlers"
	_appinfoRepositories "github.com/Rayato159/kawaii-shop/modules/appinfo/repositories"
	_appinfoUsecases "github.com/Rayato159/kawaii-shop/modules/appinfo/usecases"

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
	AppinfoModule()
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
	f.router.Get("/", _monitorHandlers.MonitorHandler(f.server.cfg).HealthCheck)
}

func (f *ModuleFactory) UsersModule() {
	repository := _usersRepositories.UsersRepository(f.server.Db())
	usecase := _usersUsecases.UsersUsecase(repository, f.server.cfg)
	handler := _usersHandlers.UsersHandler(f.server.cfg, usecase)

	router := f.router.Group("/users")

	router.Post("/signup", f.middleware.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", f.middleware.ApiKeyAuth(), handler.SignIn)
	router.Post("/signout", f.middleware.ApiKeyAuth(), handler.SignOut)
	router.Post("/refresh", f.middleware.ApiKeyAuth(), handler.RefreshPassport)

	router.Get("/secret", f.middleware.JwtAuth(), f.middleware.Authorize(2), handler.GenerateAdminToken)
	router.Get("/:user_id", f.middleware.JwtAuth(), f.middleware.ParamsCheck(), handler.GetProfile)
}

func (f *ModuleFactory) AppinfoModule() {
	repository := _appinfoRepositories.AppinfoRepository(f.server.Db())
	usecase := _appinfoUsecases.AppinfoUsecase(repository)
	handler := _appinfoHandlers.AppinfoHandler(f.server.cfg, usecase)

	router := f.router.Group("/appinfo")
	router.Get("/category", handler.FindCategory)
	router.Get("/apikey", f.middleware.JwtAuth(), f.middleware.Authorize(2), handler.GenerateApiKey)
}
