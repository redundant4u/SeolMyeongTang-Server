package session

import "github.com/labstack/echo/v4"

func registerRoutes(e *echo.Echo, h *handler) {
	e.GET("/session", h.getSessions)
	e.POST("/session", h.createSession)
	e.DELETE("/session", h.deleteSession)
}
