package httpobs

import (
	"seolmyeong-tang-server/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

const (
	echoRequestIDKey = "request_id"
	echoClientIDKey  = "client_id"
	echoTraceIDKey   = "trace_id"
)

func RequestID(c echo.Context) string {
	return fromEcho(c, echoRequestIDKey)
}

func ClientID(c echo.Context) string {
	return fromEcho(c, echoClientIDKey)
}

func TraceID(c echo.Context) string {
	return fromEcho(c, echoTraceIDKey)
}

func fromEcho(c echo.Context, key string) string {
	if c == nil {
		return ""
	}
	value, ok := c.Get(key).(string)
	if !ok {
		return ""
	}
	return value
}

// storeInEcho writes the correlation values into both the Echo context (for
// handlers) and the request context (for downstream callers using
// logger.*FromContext).
func storeInEcho(c echo.Context, requestID, clientID, traceID string) {
	c.Set(echoRequestIDKey, requestID)
	c.Set(echoClientIDKey, clientID)
	c.Set(echoTraceIDKey, traceID)
	req := c.Request()
	c.SetRequest(req.WithContext(logger.WithCorrelation(req.Context(), requestID, clientID, traceID)))
}
