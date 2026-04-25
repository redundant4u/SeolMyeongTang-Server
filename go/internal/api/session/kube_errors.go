package session

import (
	"context"
	"fmt"
	"log/slog"

	"seolmyeong-tang-server/internal/pkg/logger"
	"seolmyeong-tang-server/internal/pkg/metrics"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

type lifecycle int

const (
	lifecycleNone lifecycle = iota
	lifecyclePodCreate
	lifecyclePodDelete
)

// kubeErrorCode extracts a stable error_code label from a Kubernetes client
// error. Unknown shapes collapse to "api_error" so the label space stays small.
func kubeErrorCode(err error) string {
	if err == nil {
		return ""
	}
	if statusErr, ok := err.(kerrors.APIStatus); ok {
		if reason := string(statusErr.Status().Reason); reason != "" {
			return reason
		}
	}
	return "api_error"
}

// unwrapKubeStatus rewrites kerrors.APIStatus errors into a stable
// "code=… reason=… message=…" form so callers can return a consistent shape
// regardless of whether the raw error originated from the kube client.
func unwrapKubeStatus(err error) error {
	if err == nil {
		return nil
	}
	statusErr, ok := err.(kerrors.APIStatus)
	if !ok {
		return err
	}
	st := statusErr.Status()
	return fmt.Errorf("code=%d reason=%s message=%s", st.Code, st.Reason, st.Message)
}

// recordKubeFailure emits the three signals every kube operation failure needs:
// a metric increment (generic + lifecycle-specific), a structured log event, and
// a normalized error wrapper that callers can return directly.
func recordKubeFailure(ctx context.Context, op string, lc lifecycle, err error, attrs ...slog.Attr) error {
	code := kubeErrorCode(err)
	metrics.RecordKubeAPIFailure(op, code)
	switch lc {
	case lifecyclePodCreate:
		metrics.RecordPodCreateFailure(code)
	case lifecyclePodDelete:
		metrics.RecordPodDeleteFailure(code)
	}

	enriched := append([]slog.Attr{
		slog.String("kube_operation", op),
		slog.String("kube_error_class", code),
	}, attrs...)
	logger.ErrorEvent(ctx, "kube_api_failed", "Kubernetes API call failed", err, enriched...)

	return unwrapKubeStatus(err)
}
