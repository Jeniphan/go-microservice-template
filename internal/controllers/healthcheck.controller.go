package controllers

import (
	"net/http"
	"order-v2-microservice/internal/common/loggers"
	"order-v2-microservice/internal/services"

	"github.com/labstack/echo/v5"
)

var appLog *loggers.Logger

type HealthCheckHandler interface {
	HealthCheck(c *echo.Context) error
}

type HealthCheckController struct {
	HealthCheckService services.HealthCheckProvider
}

func NewHealthCheckController(healthCheckService services.HealthCheckProvider) HealthCheckHandler {
	appLog = loggers.NewLogger("HealthCheckController")
	return &HealthCheckController{
		HealthCheckService: healthCheckService,
	}
}

func (healthCheck *HealthCheckController) HealthCheck(c *echo.Context) error {
	appLog.Info(c, "HealthCheck Ctl")

	response, err := healthCheck.HealthCheckService.HealthCheck(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}
