package post

import (
	"seolmyeong-tang-server/internal/config"

	"github.com/labstack/echo/v4"
)

func denyPostMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Env.APP_ENV == "production" {
				return echo.ErrForbidden
			}
			return next(c)
		}
	}
}
