package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"seolmyeong-tang-server/internal/api/session"
	"seolmyeong-tang-server/internal/config"
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

func TestRouterWiringAndCORS(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile(".env", []byte("APP_ENV=production\nKUBE_SESSION_NAMESPACE=sessions\nCF_TUNNEL_ID=test\nAWS_ACCESS_KEY=test\nAWS_SECRET_KEY=test\nDYNAMODB_TABLE=posts\n"), 0600); err != nil {
		t.Fatal(err)
	}
	old := config.Env
	defer func() { config.Env = old }()
	config.InitEnv()
	e := New(nil, &session.Kube{})
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{"GET", "/metrics", 200}, {"POST", "/session/client-id", 201}, {"GET", "/session", 400}, {"GET", "/post", 403}, {"GET", "/missing", 404},
	} {
		w := httptest.NewRecorder()
		e.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		if w.Code != tc.status {
			t.Fatalf("%s %s: %d", tc.method, tc.path, w.Code)
		}
	}
	r := httptest.NewRequest("OPTIONS", "/session", nil)
	r.Header.Set("Origin", "https://redundant4u.com")
	r.Header.Set("Access-Control-Request-Method", "POST")
	r.Header.Set("Access-Control-Request-Headers", "X-Client-Id")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)
	if w.Code != 204 || w.Header().Get("Access-Control-Allow-Origin") != "https://redundant4u.com" || !strings.Contains(w.Header().Get("Access-Control-Allow-Headers"), "X-Client-Id") {
		t.Fatal("CORS preflight changed")
	}
}
