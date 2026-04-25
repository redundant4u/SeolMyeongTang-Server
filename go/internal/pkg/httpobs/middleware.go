package httpobs

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"seolmyeong-tang-server/internal/pkg/logger"
	"seolmyeong-tang-server/internal/pkg/metrics"

	"github.com/labstack/echo/v4"
)

const (
	HeaderRequestID = "X-Request-Id"
	HeaderTraceID   = "X-Trace-Id"
	HeaderClientID  = "X-Client-Id"

	maxCorrelationIDLength = 128
	unknownRoutePath       = "__unknown__"
)

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			started := time.Now()
			req := c.Request()
			requestID := clampID(req.Header.Get(HeaderRequestID))
			if requestID == "" {
				requestID = newID()
			}
			clientID := clampID(req.Header.Get(HeaderClientID))
			traceID := clampID(req.Header.Get(HeaderTraceID))

			storeInEcho(c, requestID, clientID, traceID)
			c.Response().Header().Set(HeaderRequestID, requestID)

			err := next(c)

			status := c.Response().Status
			if status == 0 {
				status = http.StatusOK
			}
			latency := time.Since(started)
			path := c.Path()
			if path == "" {
				path = unknownRoutePath
			}
			outcome := outcomeForStatus(status)

			metrics.RecordHTTPRequest(req.Method, path, strconv.Itoa(status), outcome, latency)
			logger.InfoEvent(
				c.Request().Context(),
				"http_request_completed",
				"HTTP request completed",
				slog.String("method", req.Method),
				slog.String("path", path),
				slog.Int("status", status),
				slog.Int64("latency_ms", latency.Milliseconds()),
				slog.String("ip", c.RealIP()),
				slog.String("outcome", outcome),
			)

			return err
		}
	}
}

func outcomeForStatus(status int) string {
	switch {
	case status >= 500:
		return "server_error"
	case status >= 400:
		return "client_error"
	default:
		return "success"
	}
}

func newID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(bytes[:])
}

func clampID(value string) string {
	if len(value) > maxCorrelationIDLength {
		return value[:maxCorrelationIDLength]
	}
	return value
}
