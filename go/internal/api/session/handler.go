package session

import (
	"log/slog"

	"seolmyeong-tang-server/internal/pkg/logger"
	"seolmyeong-tang-server/internal/pkg/metrics"
	"seolmyeong-tang-server/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	kube *Kube
}

func NewHandler(kube *Kube) *handler {
	return &handler{kube: kube}
}

func (h *handler) getSessions(c echo.Context) error {
	ctx := c.Request().Context()

	clientId := c.Request().Header.Get("X-Client-Id")
	if clientId == "" {
		logger.WarnEvent(ctx, "client_id_missing", "getSessions missing X-Client-Id header")
		return response.BadRequest(c)
	}

	pods, err := h.kube.getSessions(ctx, clientId)
	if err != nil {
		logger.ErrorEvent(ctx, "session_list_failed", "failed to list sessions", err)
		return response.BadRequest(c)
	}

	res, err := toGetSessionsResponse(pods)
	if err != nil {
		logger.ErrorEvent(ctx, "session_list_response_failed", "failed to map getSessions response", err)
		return response.BadRequest(c)
	}

	return response.OK(c, res)
}

func (h *handler) createSession(c echo.Context) error {
	ctx := c.Request().Context()

	clientId := c.Request().Header.Get("X-Client-Id")
	if clientId == "" {
		logger.WarnEvent(ctx, "client_id_missing", "createSession missing X-Client-Id header")
		return response.BadRequest(c)
	}

	var req createPodRequest
	if err := c.Bind(&req); err != nil {
		logger.WarnEvent(ctx, "session_create_invalid_body", "createSession invalid request body", slog.String("error", err.Error()))
		return response.BadRequest(c)
	}

	if err := c.Validate(&req); err != nil {
		logger.WarnEvent(ctx, "session_create_validation_failed", "createSession validation failed", slog.String("error", err.Error()))
		return response.BadRequest(c)
	}

	pods, err := h.kube.getSessions(ctx, clientId)
	if err != nil {
		logger.ErrorEvent(ctx, "session_list_failed", "failed to list sessions before create", err)
		return response.BadRequest(c)
	}

	if len(pods) >= 4 {
		metrics.RecordPodCreateFailure("capacity_limit")
		logger.WarnEvent(ctx, "session_create_rejected", "Session limit reached",
			slog.String("error_code", "capacity_limit"),
		)
		return response.BadRequest(c)
	}

	sessionId := h.kube.secureRandomString(8)
	metrics.RecordPodCreateRequest(req.Image)
	logger.InfoEvent(ctx, "session_create_requested", "Session creation requested",
		slog.String("session_id", sessionId),
		slog.String("requested_image", req.Image),
	)

	info := createPod{
		name:        req.Name,
		image:       req.Image,
		clientId:    clientId,
		sessionId:   sessionId,
		description: req.Description,
	}

	pod, err := h.kube.createSession(ctx, info)
	if err != nil {
		logger.ErrorEvent(ctx, "session_create_failed", "Session creation failed", err,
			slog.String("session_id", sessionId),
		)
		return response.BadRequest(c)
	}
	logger.InfoEvent(ctx, "session_create_succeeded", "Session creation succeeded",
		slog.String("session_id", sessionId),
		slog.String("pod_name", pod.Name),
	)

	res, err := toCreateSessionResponse(pod, sessionId)
	if err != nil {
		logger.ErrorEvent(ctx, "session_create_response_failed", "failed to convert createSession response", err)
		return response.BadRequest(c)
	}

	return response.Created(c, res)
}

func (h *handler) deleteSession(c echo.Context) error {
	ctx := c.Request().Context()

	clientId := c.Request().Header.Get("X-Client-Id")
	if clientId == "" {
		logger.WarnEvent(ctx, "client_id_missing", "deleteSession missing X-Client-Id header")
		return response.BadRequest(c)
	}

	var req deletePodRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c)
	}

	if err := c.Validate(&req); err != nil {
		logger.WarnEvent(ctx, "session_delete_validation_failed", "deleteSession validation failed", slog.String("error", err.Error()))
		return response.BadRequest(c)
	}

	info := deletePod{
		clientId:  clientId,
		sessionId: req.SessionId,
	}

	metrics.RecordPodDeleteRequest()
	logger.InfoEvent(ctx, "session_delete_requested", "Session deletion requested",
		slog.String("session_id", req.SessionId),
	)

	if err := h.kube.deleteSession(ctx, info); err != nil {
		logger.ErrorEvent(ctx, "session_delete_failed", "Session deletion failed", err,
			slog.String("session_id", req.SessionId),
		)
		return response.BadRequest(c)
	}
	logger.InfoEvent(ctx, "session_delete_succeeded", "Session deletion succeeded",
		slog.String("session_id", req.SessionId),
	)

	return response.NoContent(c)
}

func (h *handler) createClientId(c echo.Context) error {
	ctx := c.Request().Context()
	clientId := h.kube.secureRandomString(8)
	logger.InfoEvent(ctx, "client_id_generated", "Client identifier generated",
		slog.String("client_id", clientId),
		slog.String("client_display_id", logger.ClientDisplayID(clientId)),
	)

	res := createClientIdResponse{
		ClientId: clientId,
	}

	return response.Created(c, res)
}
