package loggers

import (
	"encoding/json"
	"fmt"
	"log"
	"order-v2-microservice/internal/models"
	"time"

	"github.com/labstack/echo/v5"
)

type Logger struct {
	Service string
}

func (rl *Logger) Info(c *echo.Context, message ...interface{}) {
	reqId := c.Request().Header.Get(echo.HeaderXRequestID)

	logging := models.CommonLog{
		TimeStamp:   time.Now(),
		Service:     rl.Service,
		RequestId:   reqId,
		LoggingType: "ApiLogs",
		LogLevel:    "Info",
		Message:     fmt.Sprint(message...),
	}

	logMessage := loggingToJson(logging)
	log.Println(logMessage)
}

func (rl *Logger) Error(c *echo.Context, message ...interface{}) {
	reqId := c.Request().Header.Get(echo.HeaderXRequestID)

	logging := models.CommonLog{
		TimeStamp:   time.Now(),
		Service:     rl.Service,
		RequestId:   reqId,
		LoggingType: "ApiLogs",
		LogLevel:    "Error",
		Message:     fmt.Sprint(message...),
	}

	logMessage := loggingToJson(logging)
	log.Println(logMessage)
}

func loggingToJson(logging models.CommonLog) string {
	b, err := json.Marshal(logging)
	if err != nil {
		log.Printf("Converse logging to json err: %s, service: %s, request_id %s \n", err, logging.Service, logging.RequestId)
		return ""
	}

	return string(b)
}

func NewLogger(service string) *Logger {
	return &Logger{
		Service: service,
	}
}
