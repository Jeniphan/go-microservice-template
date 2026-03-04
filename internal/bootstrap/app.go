package app

import (
	"order-v2-microservice/configs"
	"order-v2-microservice/internal/controllers"
	"order-v2-microservice/internal/middlewares"
	"order-v2-microservice/internal/services"

	"gorm.io/gorm"
)

type Handler struct {
	*Middlewares
	*Controllers
}

type Services struct {
	HealthCheckService services.HealthCheckProvider
}

type Controllers struct {
	HealthCheckCtrl controllers.HealthCheckHandler
}

type Middlewares struct {
	AppMdw *middlewares.AppMiddleware
}

type AppsProvider interface {
	CreateAppCtrl()
	CreateAppMdw()
	CreateService()
}

type Apps struct {
	DB          *gorm.DB
	Services    *Services
	Controllers *Controllers
	Middlewares *Middlewares
}

func (a *Apps) CreateService() {
	healthCheck := services.NewHealthCheckService(a.DB)

	a.Services = &Services{
		HealthCheckService: healthCheck,
	}
}

func (a *Apps) CreateAppMdw() {
	mdw := middlewares.NewApplicationMiddleware()
	a.Middlewares = &Middlewares{
		AppMdw: mdw,
	}
}

func (a *Apps) CreateAppCtrl() {
	HealthCheckCtrl := controllers.NewHealthCheckController(a.Services.HealthCheckService)

	a.Controllers = &Controllers{
		HealthCheckCtrl: HealthCheckCtrl,
	}
}

func Bootstrap() *Handler {
	configs.InitConfigEnv()
	db := configs.NewDatabaseConfig()

	app := &Apps{
		DB: db,
	}

	app.CreateService()
	app.CreateAppCtrl()
	app.CreateAppMdw()

	handler := &Handler{
		Middlewares: app.Middlewares,
		Controllers: app.Controllers,
	}

	return handler
}
