package post

import (
	"seolmyeong-tang-server/internal/pkg/logger"
	"seolmyeong-tang-server/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	repo *repository
}

func newHandler(r *repository) *handler {
	return &handler{repo: r}
}

func (h *handler) getPosts(c echo.Context) error {
	ctx := c.Request().Context()

	posts, err := h.repo.getPosts(ctx)
	if err != nil {
		logger.Error("failed to get posts", err)
		return response.BadRequest(c)
	}

	res := toGetPostsResponse(posts)

	return response.OK(c, res)
}

func (h *handler) getPost(c echo.Context) error {
	ctx := c.Request().Context()

	postId := c.Param("postId")
	if postId == "" {
		logger.Error("postId is required", nil)
		return response.BadRequest(c)
	}

	post, err := h.repo.getPost(ctx, postId)
	if err != nil {
		logger.Error("failed to get post", err)
		return response.BadRequest(c)
	}

	res := toGetPostResponse(post)

	return response.OK(c, res)
}
