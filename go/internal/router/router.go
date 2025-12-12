package router

import (
	"seolmyeong-tang-server/internal/api/post"
	"seolmyeong-tang-server/internal/api/session"
	"seolmyeong-tang-server/internal/pkg/logger"
	"seolmyeong-tang-server/internal/pkg/validator"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(ddb *dynamodb.Client) *echo.Echo {
	e := echo.New()

	e.Validator = validator.New()

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

	session.Init(e)
	post.Init(e, ddb)

	logger.Info("Router init")

	return e
}
