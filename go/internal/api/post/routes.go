package post

import (
	"seolmyeong-tang-server/internal/config"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, ddb *dynamodb.Client) error {
	repo := newRepository(ddb, config.Env.DYNAMODB_TABLE)
	handler := newHandler(repo)

	registerRoutes(e, handler)

	return nil
}

func registerRoutes(e *echo.Echo, h *handler) {
	postGroup := e.Group("/post")

	postGroup.Use(denyPostMiddleware())

	postGroup.GET("/post/:postId", h.getPost)
	postGroup.GET("/post", h.getPosts)
}
