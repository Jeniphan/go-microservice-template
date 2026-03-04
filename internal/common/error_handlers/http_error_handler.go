package errorhandler

import (
	"errors"
	"net/http"
	"order-v2-microservice/internal/common/loggers"

	"github.com/labstack/echo/v5"
)

var log *loggers.Logger

func init() {
	log = loggers.NewLogger("FilterHTTPErrorHandler")
}

func FilterHTTPErrorHandler(c *echo.Context, err error) {
	// Log the error
	//c.Logger().Error(err)
	log.Info(c, err.Error())

	// Check if the error is an *echo.HTTPError
	var he *echo.HTTPError
	if errors.As(err, &he) {
		// Send custom response
		c.JSON(he.Code, map[string]interface{}{
			"status":  "error",
			"message": he.Message,
		})
		return
	}

	// For non-HTTP errors, send a generic 500 response
	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"status":  "error",
		"message": "Internal Server Error",
	})
}
