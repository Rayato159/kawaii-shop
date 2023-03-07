package servers

import (
	_middlewareHandlers "github.com/Rayato159/kawaii-shop/modules/middlewares/handlers"
	_middlewareRepositories "github.com/Rayato159/kawaii-shop/modules/middlewares/repositories"
	_middlewareUsecases "github.com/Rayato159/kawaii-shop/modules/middlewares/usecases"

	_monitorHandlers "github.com/Rayato159/kawaii-shop/modules/monitor/handlers"

	_filesHandlers "github.com/Rayato159/kawaii-shop/modules/files/handlers"
	_filesRepositories "github.com/Rayato159/kawaii-shop/modules/files/repositories"
	_filesUsecases "github.com/Rayato159/kawaii-shop/modules/files/usecases"

	_usersHandlers "github.com/Rayato159/kawaii-shop/modules/users/handlers"
	_usersRepositories "github.com/Rayato159/kawaii-shop/modules/users/repositories"
	_usersUsecases "github.com/Rayato159/kawaii-shop/modules/users/usecases"

	_appinfoHandlers "github.com/Rayato159/kawaii-shop/modules/appinfo/handlers"
	_appinfoRepositories "github.com/Rayato159/kawaii-shop/modules/appinfo/repositories"
	_appinfoUsecases "github.com/Rayato159/kawaii-shop/modules/appinfo/usecases"

	_productsHandlers "github.com/Rayato159/kawaii-shop/modules/products/handlers"
	_productsRepositories "github.com/Rayato159/kawaii-shop/modules/products/repositories"
	_productsUsecases "github.com/Rayato159/kawaii-shop/modules/products/usecases"

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
	FilesModule()
	UsersModule()
	AppinfoModule()
	ProductsModule()
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

func (f *ModuleFactory) FilesModule() {
	repository := _filesRepositories.FilesRepository(f.server.Db())
	usecase := _filesUsecases.FilesUsecase(repository)
	handler := _filesHandlers.FilesHandler(f.server.cfg, usecase)

	router := f.router.Group("/files")

	router.Post("/", f.middleware.JwtAuth(), handler.UploadFiles)
}

func (f *ModuleFactory) UsersModule() {
	repository := _usersRepositories.UsersRepository(f.server.Db())
	usecase := _usersUsecases.UsersUsecase(repository, f.server.cfg)
	handler := _usersHandlers.UsersHandler(f.server.cfg, usecase)

	router := f.router.Group("/users")

	router.Post("/signup", f.middleware.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/admin", f.middleware.JwtAuth(), f.middleware.Authorize(2), handler.AddAdmin)
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

	router.Post("/categories", f.middleware.JwtAuth(), f.middleware.Authorize(2), handler.AddCategory)

	router.Get("/categories", f.middleware.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", f.middleware.JwtAuth(), f.middleware.Authorize(2), handler.GenerateApiKey)

	router.Delete("/categories", f.middleware.JwtAuth(), f.middleware.Authorize(2), handler.RemoveCategory)
}

func (f *ModuleFactory) ProductsModule() {
	repository := _productsRepositories.ProductsRepository(f.server.Db())
	usecase := _productsUsecases.ProductsUsecase(repository)
	handler := _productsHandlers.ProductsHandler(f.server.cfg, usecase)

	router := f.router.Group("/products")

	router.Get("/", f.middleware.ApiKeyAuth(), handler.FindProduct)
	router.Get("/:product_id", f.middleware.ApiKeyAuth(), handler.FindOneProduct)

	router.Post("/", f.middleware.JwtAuth(), f.middleware.Authorize(2), handler.AddProduct)
}
