package session

import (
	"fmt"
	"seolmyeong-tang-server/internal/config"
	"seolmyeong-tang-server/internal/pkg/k8s"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) error {
	kubeClient, err := k8s.NewClient(config.Env.KUBE_CONFIG)
	if err != nil {
		return fmt.Errorf("failed to initialize kube client: %w", err)
	}

	kube := newKube(kubeClient, config.Env.KUBE_SESSION_NAMESPACE)
	handler := NewHandler(kube)

	registerRoutes(e, handler)
	return nil
}
