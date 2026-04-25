package router

import (
	"context"

	"seolmyeong-tang-server/internal/api/post"
	"seolmyeong-tang-server/internal/api/session"
	"seolmyeong-tang-server/internal/pkg/httpobs"
	"seolmyeong-tang-server/internal/pkg/logger"
	"seolmyeong-tang-server/internal/pkg/metrics"
	"seolmyeong-tang-server/internal/pkg/validator"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(ddb *dynamodb.Client, kube *session.Kube) *echo.Echo {
	e := echo.New()

	e.Validator = validator.New()

	e.Use(middleware.Recover())
	e.Use(httpobs.Middleware())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://redundant4u.com"},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderContentType,
			echo.HeaderAuthorization,
			httpobs.HeaderRequestID,
			httpobs.HeaderTraceID,
			"X-Client-Id",
		},
	}))

	e.GET("/metrics", echo.WrapHandler(metrics.Handler()))

	session.Init(e, kube)
	post.Init(e, ddb)

	logger.InfoEvent(context.Background(), "router_initialized", "Router initialized")

	return e
}
