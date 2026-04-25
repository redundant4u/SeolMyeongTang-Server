package httpobs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"seolmyeong-tang-server/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

func TestMiddlewareAddsCorrelationContextAndHeader(t *testing.T) {
	e := echo.New()
	e.Use(Middleware())
	e.GET("/health", func(c echo.Context) error {
		if RequestID(c) == "" {
			t.Fatal("request id was not stored in echo context")
		}
		if ClientID(c) != "client-1" {
			t.Fatalf("unexpected client id: %q", ClientID(c))
		}
		if logger.RequestIDFromContext(c.Request().Context()) == "" {
			t.Fatal("request id was not stored in request context")
		}
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set(HeaderClientID, "client-1")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if rec.Header().Get(HeaderRequestID) == "" {
		t.Fatal("request id response header was not set")
	}
}
