package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"order-v2-microservice/internal/common/loggers"
	"order-v2-microservice/internal/models/responses"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
)

var appLog *loggers.Logger

type App interface {
	AuthenticatedToken(next echo.HandlerFunc) echo.HandlerFunc
}

type AppMiddleware struct{}

func NewApplicationMiddleware() *AppMiddleware {
	appLog = loggers.NewLogger("AppMiddleware")
	return &AppMiddleware{}
}

func (a *AppMiddleware) AuthenticatedToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		requestID := c.Request().Header.Get(echo.HeaderXRequestID)
		appLog.Info(c, "Authentication middleware triggered", "requestId", requestID)

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) < 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token or Bearer token")
		}

		idToken := parts[1]

		// Verify token by calling API
		requestBody := map[string]string{"token": idToken}
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			appLog.Error(c, "Failed to marshal request body", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to process token verification")
		}

		req, err := http.NewRequest(os.Getenv("METHOD_APPLICATION_VERIFY"), os.Getenv("ENDPOINT_APPLICATION_MS")+os.Getenv("URL_APPLICATION_VERIFY"), bytes.NewBuffer(jsonBody))
		if err != nil {
			appLog.Error(c, "Failed to create request", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create verification request")
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "*/*")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			appLog.Error(c, "Failed to call verify token API", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to verify token")
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			appLog.Error(c, "Failed to read response body", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to read verification response")
		}

		var verifyResponse responses.ApplicationVerifyResponse
		if err := json.Unmarshal(body, &verifyResponse); err != nil {
			appLog.Error(c, "Failed to unmarshal response", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse verification response")
		}

		if resp.StatusCode != http.StatusOK || verifyResponse.Status != 200 {
			appLog.Info(c, "Token verification failed", "status", verifyResponse.Status, "message", verifyResponse.Message)
			return echo.NewHTTPError(http.StatusUnauthorized, "token verification failed: "+verifyResponse.Message)
		}

		appLog.Info(c, "Token verified successfully", "result", verifyResponse.Result)
		c.Request().Header.Set("AppID", verifyResponse.Result)

		return next(c)
	}
}
