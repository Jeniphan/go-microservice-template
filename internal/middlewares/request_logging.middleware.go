package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"order-v2-microservice/internal/models"
	"time"

	"github.com/labstack/echo/v5"
)

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func RequestLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// Log header details
		headers, headerToJsonErr := flattenMapHeader(c.Request().Header)
		if headerToJsonErr != nil {
			log.Printf("Failed to read request header: %v", headerToJsonErr)
		}

		// Log request details
		var reqBody []byte
		if c.Request().Body != nil { // Read
			reqBody, _ = io.ReadAll(c.Request().Body)
		}
		c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

		start := time.Now()

		// Capture response body (requires a custom writer)
		resBody := new(bytes.Buffer)
		origWriter := c.Response()
		mw := io.MultiWriter(origWriter, resBody)
		writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: origWriter}
		c.SetResponse(writer)

		if err := next(c); err != nil {
			//c.Logger().Error(err) // Log the error
			return err // Return the error to be handled once by the global handler
		}

		// Get status after handler runs via type assertion to *echo.Response
		var status int
		if echoResp, ok := origWriter.(*echo.Response); ok {
			status = echoResp.Status
		}
		latency := time.Since(start)

		logging := models.RequestLog{
			Timestamp:   time.Now(),
			Service:     "RequestLoggingMiddleware",
			RequestId:   c.Request().Header.Get(echo.HeaderXRequestID),
			LoggingType: "RequestLogs",
			Hostname:    c.Request().Host,
			Uri:         c.Request().RequestURI,
			Method:      c.Request().Method,
			Path:        c.Path(),
			Start:       start,
			End:         time.Now(),
			Headers:     headers,
			Request:     string(reqBody),
			Response:    string(resBody.Bytes()),
			Status:      status,
			Latency:     latency,
		}

		b, logToJsonError := json.Marshal(logging)
		if logToJsonError != nil {
			log.Fatalf("Failed to format request response logging : %v", logToJsonError)
		}
		log.Println(string(b))
		return nil
	}
}

func flattenMapHeader(input http.Header) (string, error) {
	flattened := make(map[string]string)

	// Flatten the map by taking the first value of each slice
	for key, values := range input {
		if len(values) > 0 {
			flattened[key] = values[0]
		}
	}

	// Convert the flattened map to a JSON string
	jsonData, err := json.Marshal(flattened)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
