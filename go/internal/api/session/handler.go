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
	clientId := c.Request().Header.Get("X-Client-Id")
	if clientId == "" {
		logger.Error("getSessions client-id header is not set", nil)
	}

	ctx := c.Request().Context()
	pods, err := h.kube.getSessions(ctx, clientId)
	if err != nil {
		logger.Error("failed to get pods", err)
		return response.BadRequest(c)
	}

	res, err := toGetSessionsResponse(pods)
	if err != nil {
		logger.Error("failed to map getSessions response", err)
		return response.BadRequest(c)
	}

	return response.OK(c, res)
}

func (h *handler) createSession(c echo.Context) error {
	clientId := c.Request().Header.Get("X-Client-Id")
	if clientId == "" {
		logger.Error("createSession X-Client-Id header is not set", nil)
		return response.BadRequest(c)
	}

	var req createPodRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("createSession invalid request body", nil)
		return response.BadRequest(c)
	}

	if req.Name == "" {
		logger.Error("createSession body value is invalid", nil)
		return response.BadRequest(c)
	}

	ctx := c.Request().Context()

	pods, err := h.kube.getSessions(ctx, clientId)
	if err != nil {
		logger.Error("failed to get pods", err)
		return response.BadRequest(c)
	}

	if len(pods) >= 4 {
		logger.Error("session limit reached", nil)
		return response.BadRequest(c)
	}

	sessionId := h.kube.secureRandomString(8)

	info := createPod{
		name:      req.Name,
		clientId:  clientId,
		sessionId: sessionId,
	}

	pod, err := h.kube.createSession(ctx, info)
	if err != nil {
		logger.Error("failed to create pods", err)
		return response.BadRequest(c)
	}

	logger.Info("pod is created: %s", sessionId)

	res, err := toCreateSessionResponse(pod, sessionId)
	if err != nil {
		logger.Error("failed to convert createSession response", err)
		return response.BadRequest(c)
	}

	return response.Created(c, res)
}

func (h *handler) deleteSession(c echo.Context) error {
	clientId := c.Request().Header.Get("X-Client-Id")
	if clientId == "" {
		logger.Error("deleteSession X-Client-Id header is not set", nil)
		return response.BadRequest(c)
	}

	var req deletePodRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c)
	}

	ctx := c.Request().Context()

	if req.SessionId == "" {
		logger.Error("sessionId is required", nil)
		return response.BadRequest(c)
	}

	info := deletePod{
		clientId:  clientId,
		sessionId: req.SessionId,
	}

	err := h.kube.deleteSession(ctx, info)
	if err != nil {
		logger.Error("failed to delete session", err)
		return response.BadRequest(c)
	}

	logger.Info("pod is deleted: %s", req.SessionId)

	return response.NoContent(c)
}

func (h *handler) createClientId(c echo.Context) error {
	clientId := h.kube.secureRandomString(8)
	logger.Info("generated client id: %s", clientId)

	res := createClientIdResponse{
		ClientId: clientId,
	}

	return response.Created(c, res)
}
