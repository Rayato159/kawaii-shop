package servers

import (
	_middlewareHandlers "github.com/Rayato159/kawaii-shop/modules/middlewares/handlers"
	_middlewareRepositories "github.com/Rayato159/kawaii-shop/modules/middlewares/repositories"
	_middlewareUsecases "github.com/Rayato159/kawaii-shop/modules/middlewares/usecases"

	_monitorHandlers "github.com/Rayato159/kawaii-shop/modules/monitor/handlers"

	_filesHandlers "github.com/Rayato159/kawaii-shop/modules/files/handlers"
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

	_ordersHandlers "github.com/Rayato159/kawaii-shop/modules/orders/handlers"
	_ordersRepositories "github.com/Rayato159/kawaii-shop/modules/orders/repositories"
	_ordersUsecases "github.com/Rayato159/kawaii-shop/modules/orders/usecases"

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
	OrdersModule()
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
	usecase := _filesUsecases.FilesUsecase(f.server.cfg)
	handler := _filesHandlers.FilesHandler(f.server.cfg, usecase)

	router := f.router.Group("/files")

	router.Post("/", f.middleware.JwtAuth(), handler.UploadFiles)

	router.Patch("/", f.middleware.JwtAuth(), handler.DeleteFile)
}

func (f *ModuleFactory) UsersModule() {
	repository := _usersRepositories.UsersRepository(f.server.db)
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
	repository := _appinfoRepositories.AppinfoRepository(f.server.db)
	usecase := _appinfoUsecases.AppinfoUsecase(repository)
	handler := _appinfoHandlers.AppinfoHandler(f.server.cfg, usecase)

	router := f.router.Group("/appinfo")

	router.Post("/categories", f.middleware.JwtAuth(), f.middleware.Authorize(2), handler.AddCategory)

	router.Get("/categories", f.middleware.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", f.middleware.JwtAuth(), f.middleware.Authorize(2), handler.GenerateApiKey)

	router.Delete("/categories", f.middleware.JwtAuth(), f.middleware.Authorize(2), handler.RemoveCategory)
}

func (f *ModuleFactory) ProductsModule() {
	// File Module
	filesUsecase := _filesUsecases.FilesUsecase(f.server.cfg)

	productsRepository := _productsRepositories.ProductsRepository(f.server.db, f.server.cfg, filesUsecase)
	productsUsecase := _productsUsecases.ProductsUsecase(productsRepository)
	productsHandler := _productsHandlers.ProductsHandler(f.server.cfg, productsUsecase, filesUsecase)

	router := f.router.Group("/products")

	router.Get("/", f.middleware.ApiKeyAuth(), productsHandler.FindProduct)
	router.Get("/:product_id", f.middleware.ApiKeyAuth(), productsHandler.FindOneProduct)

	router.Post("/", f.middleware.JwtAuth(), f.middleware.Authorize(2), productsHandler.AddProduct)

	router.Patch("/:product_id", f.middleware.JwtAuth(), f.middleware.Authorize(2), productsHandler.UpdateProduct)

	router.Delete("/:product_id", f.middleware.JwtAuth(), f.middleware.Authorize(2), productsHandler.DeleteProduct)
}

func (f *ModuleFactory) OrdersModule() {
	filesUsecase := _filesUsecases.FilesUsecase(f.server.cfg)
	productsRepository := _productsRepositories.ProductsRepository(f.server.db, f.server.cfg, filesUsecase)

	ordersRepository := _ordersRepositories.OrdersRepository(f.server.db)
	ordersUsecase := _ordersUsecases.OrdersUsecase(ordersRepository, productsRepository)
	ordersHandler := _ordersHandlers.OrdersHandler(f.server.cfg, ordersUsecase)

	router := f.router.Group("/orders")

	router.Get("/", f.middleware.JwtAuth(), f.middleware.Authorize(2), ordersHandler.FindOrder)
	router.Get("/:order_id", f.middleware.JwtAuth(), f.middleware.Authorize(2), ordersHandler.FindOneOrder)

	router.Post("/", f.middleware.JwtAuth(), ordersHandler.CreateOrder)

	router.Patch("/:order_id", f.middleware.JwtAuth(), ordersHandler.UpdateOrder)
}
