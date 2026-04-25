package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"seolmyeong-tang-server/internal/pkg/metrics"

	"github.com/labstack/echo/v4"
)

func TestMetricsEndpointIsRegistered(t *testing.T) {
	metrics.RecordHTTPRequest("GET", "/session", "200", "success", time.Millisecond)
	e := echo.New()
	e.GET("/metrics", echo.WrapHandler(metrics.Handler()))
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "smt_http_requests_total") {
		t.Fatal("metrics response does not include SMT HTTP collector")
	}
}
