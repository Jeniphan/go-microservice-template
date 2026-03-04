package services

import (
	"order-v2-microservice/internal/common/loggers"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

var appLog *loggers.Logger

type HealthCheckProvider interface {
	HealthCheck(c *echo.Context) (string, error)
}

type HealthCheckService struct {
	DB *gorm.DB
}

func NewHealthCheckService(db *gorm.DB) HealthCheckProvider {
	appLog = loggers.NewLogger("HealthCheckService")
	healthCheckService := &HealthCheckService{
		DB: db,
	}

	return healthCheckService
}

func (healthCheck *HealthCheckService) HealthCheck(c *echo.Context) (string, error) {
	appLog.Info(c, "HealthCheck")
	return "Good bye world", nil
}
