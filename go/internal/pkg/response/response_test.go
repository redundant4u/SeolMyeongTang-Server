package response

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatusAndPublicErrorBodies(t *testing.T) {
	for _, tc := range []struct {
		code int
		call func(echo.Context) error
	}{
		{200, func(c echo.Context) error { return OK(c, echo.Map{"value": "ok"}) }},
		{201, func(c echo.Context) error { return Created(c, echo.Map{"value": "ok"}) }},
		{204, NoContent}, {404, NotFound}, {400, BadRequest},
		{500, func(c echo.Context) error { return InternalError(c, errors.New("private provider details")) }},
	} {
		w := httptest.NewRecorder()
		c := echo.New().NewContext(httptest.NewRequest("GET", "/", nil), w)
		if err := tc.call(c); err != nil || w.Code != tc.code {
			t.Fatalf("status %d: %d %v", tc.code, w.Code, err)
		}
		if strings.Contains(w.Body.String(), "private") {
			t.Fatal("internal error details leaked")
		}
	}
}
