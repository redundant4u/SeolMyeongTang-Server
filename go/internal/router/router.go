package router

import (
	"seolmyeong-tang-server/internal/api/session"
	"seolmyeong-tang-server/internal/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() (*echo.Echo, error) {
	e := echo.New()

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3001"},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderContentType,
			echo.HeaderAuthorization,
			"X-Client-Id",
		},
	}))

	if err := session.Init(e); err != nil {
		return nil, err
	}

	logger.Info("Router init")

	return e, nil
}
