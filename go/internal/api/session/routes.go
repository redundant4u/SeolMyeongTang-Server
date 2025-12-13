package session

import (
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, kube *Kube) error {
	handler := NewHandler(kube)

	registerRoutes(e, handler)

	return nil
}

func registerRoutes(e *echo.Echo, h *handler) {
	e.GET("/session", h.getSessions)

	e.POST("/session", h.createSession)
	e.POST("/session/client-id", h.createClientId)

	e.DELETE("/session", h.deleteSession)
}
