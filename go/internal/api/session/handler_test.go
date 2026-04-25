package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"seolmyeong-tang-server/internal/pkg/httpobs"
	"seolmyeong-tang-server/internal/pkg/validator"

	"github.com/labstack/echo/v4"
)

func newTestEcho(h *handler) *echo.Echo {
	e := echo.New()
	e.Validator = validator.New()
	e.Use(httpobs.Middleware())
	e.POST("/session/client-id", h.createClientId)
	e.POST("/session", h.createSession)
	e.DELETE("/session", h.deleteSession)
	return e
}

func TestCreateClientIDEmitsCreatedResponse(t *testing.T) {
	e := newTestEcho(NewHandler(&Kube{}))

	req := httptest.NewRequest(http.MethodPost, "/session/client-id", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "clientId") {
		t.Fatalf("response did not include clientId: %s", rec.Body.String())
	}
	if rec.Header().Get(httpobs.HeaderRequestID) == "" {
		t.Fatal("request id response header was not set")
	}
}

func TestCreateSessionMissingClientIDReturns400(t *testing.T) {
	e := newTestEcho(NewHandler(&Kube{}))

	body := `{"name":"test","image":"ubuntu","description":"desc"}`
	req := httptest.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing X-Client-Id, got %d", rec.Code)
	}
}

func TestCreateSessionInvalidBodyReturns400(t *testing.T) {
	e := newTestEcho(NewHandler(&Kube{}))

	req := httptest.NewRequest(http.MethodPost, "/session", strings.NewReader("not json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("X-Client-Id", "test-client")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad body, got %d", rec.Code)
	}
}

func TestCreateSessionValidationFailureReturns400(t *testing.T) {
	e := newTestEcho(NewHandler(&Kube{}))

	body := `{"name":"test","image":"invalid-image","description":"desc"}`
	req := httptest.NewRequest(http.MethodPost, "/session", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("X-Client-Id", "test-client")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for validation failure, got %d", rec.Code)
	}
}

func TestDeleteSessionMissingClientIDReturns400(t *testing.T) {
	e := newTestEcho(NewHandler(&Kube{}))

	body := `{"sessionId":"abc12345"}`
	req := httptest.NewRequest(http.MethodDelete, "/session", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing X-Client-Id, got %d", rec.Code)
	}
}
