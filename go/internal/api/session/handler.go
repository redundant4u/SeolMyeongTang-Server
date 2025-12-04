package session

import (
	"seolmyeong-tang-server/internal/pkg/logger"
	"seolmyeong-tang-server/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	kube *kube
}

func NewHandler(kube *kube) *handler {
	return &handler{kube: kube}
}

func (h *handler) getSessions(c echo.Context) error {
	ctx := c.Request().Context()
	pods, err := h.kube.getSessions(ctx)
	if err != nil {
		logger.Error("failed to get pods", err)
		return response.BadRequest(c)
	}

	res := make([]getPodsResponse, 0, len(pods))
	for _, p := range pods {
		sessionId, ok := p.Labels["session-id"]
		if !ok || sessionId == "" {
			logger.Error("pod session-id is something wrong", err)
			return response.BadRequest(c)
		}

		res = append(res, getPodsResponse{
			Name:      p.Name,
			SessionId: sessionId,
		})
	}

	return response.OK(c, res)
}

func (h *handler) createSession(c echo.Context) error {
	ctx := c.Request().Context()
	pods, err := h.kube.getSessions(ctx)
	if err != nil {
		logger.Error("failed to get pods", err)
		return response.BadRequest(c)
	}

	if len(pods) >= 4 {
		logger.Error("session limit reached", nil)
		return response.BadRequest(c)
	}

	name := "vnc"
	sessionId := h.kube.secureRandomString(8)

	info := createPodRequest{
		Name:      name,
		SessionId: sessionId,
	}

	pod, err := h.kube.createSession(ctx, info)
	if err != nil {
		logger.Error("failed to create pods", err)
		return response.BadRequest(c)
	}

	res := createPodResponse{
		Name:      pod.Name,
		SessionId: sessionId,
	}

	return response.Created(c, res)
}

func (h *handler) deleteSession(c echo.Context) error {
	ctx := c.Request().Context()

	var req deletePodRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c)
	}

	if req.SessionId == "" {
		logger.Error("sessionId is required", nil)
		return response.BadRequest(c)
	}

	err := h.kube.deleteSession(ctx, req)
	if err != nil {
		logger.Error("failed to delete session", err)
		return response.BadRequest(c)
	}

	return response.NoContent(c)
}
