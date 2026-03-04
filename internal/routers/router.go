package routers

import (
	"order-v2-microservice/configs"
	app "order-v2-microservice/internal/bootstrap"
	errorhandler "order-v2-microservice/internal/common/error_handlers"
	appmiddleware "order-v2-microservice/internal/middlewares"

	"github.com/labstack/echo/v5"
)

func SetupRouter(h *app.Handler) *echo.Echo {
	e := echo.New()
	// Validation
	e.Validator = configs.NewValidator()

	// Error handler
	e.HTTPErrorHandler = errorhandler.FilterHTTPErrorHandler
	e.Use(appmiddleware.RequestLogging)

	// HealthCheck routes
	e.GET("/order/healthcheck", h.HealthCheckCtrl.HealthCheck)
	_ = e.Group("/order", h.AppMdw.AuthenticatedToken)

	return e
}
