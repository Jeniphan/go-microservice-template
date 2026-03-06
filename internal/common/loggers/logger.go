package loggers

import (
	"encoding/json"
	"fmt"
	"log"
	"order-v2-microservice/internal/models"
	"os"
	"time"

	"github.com/labstack/echo/v5"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
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

	fmt.Print(formatLog(logging))
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

	fmt.Print(formatLog(logging))
}

func formatLog(logging models.CommonLog) string {
	pid := os.Getpid()
	timestamp := logging.TimeStamp.Format("01/02/2006, 3:04:05 PM")

	var levelColor, levelStr string
	switch logging.LogLevel {
	case "Error":
		levelColor = colorRed
		levelStr = "ERROR"
	default:
		levelColor = colorGreen
		levelStr = "LOG  "
	}

	return fmt.Sprintf(
		"%s[AkiiraTech-Logger]%s %d  - %s     %s%s%s %s[%s]%s %s (req_id: %s)\n",
		colorYellow, colorReset,
		pid,
		timestamp,
		levelColor, levelStr, colorReset,
		colorYellow, logging.Service, colorReset,
		logging.Message,
		logging.RequestId,
	)
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
