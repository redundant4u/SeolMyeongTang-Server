package session

import (
	"context"
	"log/slog"
	"time"

	"seolmyeong-tang-server/internal/pkg/k8s"
	"seolmyeong-tang-server/internal/pkg/logger"
	"seolmyeong-tang-server/internal/pkg/metrics"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type gc struct {
	k8s       *k8s.Client
	namespace string
	interval  int
}

func newGC(k8s *k8s.Client, namespace string, interval int) *gc {
	return &gc{
		k8s:       k8s,
		namespace: namespace,
		interval:  interval,
	}
}

func (g *gc) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(g.interval) * time.Second)
	defer ticker.Stop()

	logger.InfoEvent(ctx, "gc_started", "Session garbage collector started",
		slog.Int("interval_seconds", g.interval),
		slog.String("namespace", g.namespace),
	)
	for {
		select {
		case <-ticker.C:
			g.cleanup(ctx)
		case <-ctx.Done():
			logger.InfoEvent(ctx, "gc_stopped", "Session garbage collector stopped",
				slog.String("namespace", g.namespace),
			)
			return
		}
	}
}

func (g *gc) cleanup(ctx context.Context) {
	listOptions := metav1.ListOptions{
		LabelSelector: "app=vnc",
		Limit:         100,
	}

	activeCount := 0

	for {
		pods, err := g.k8s.Clientset.CoreV1().Pods(g.namespace).List(ctx, listOptions)
		if err != nil {
			_ = recordKubeFailure(ctx, "gc_list_pods", lifecycleNone, err,
				slog.String("namespace", g.namespace),
			)
			return
		}

		now := time.Now().UTC()
		for _, pod := range pods.Items {
			if pod.ObjectMeta.DeletionTimestamp == nil &&
				(pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodPending) {
				activeCount++
			}

			if pod.ObjectMeta.DeletionTimestamp != nil {
				continue
			}

			expiredAtStr, ok := pod.Annotations["expired-at"]
			if !ok {
				continue
			}

			expiredAt, err := time.Parse(time.RFC3339, expiredAtStr)
			if err != nil {
				logger.ErrorEvent(ctx, "gc_expired_at_parse_failed", "failed to parse expired-at annotation", err,
					slog.String("pod_name", pod.Name),
					slog.String("namespace", g.namespace),
				)
				continue
			}

			if now.After(expiredAt) {
				g.deleteExpired(ctx, pod)
			}
		}

		if pods.Continue == "" {
			break
		}
		listOptions.Continue = pods.Continue
	}

	metrics.SetActiveSessions(activeCount)
}

func (g *gc) deleteExpired(ctx context.Context, pod corev1.Pod) {
	clientId := pod.Labels["client-id"]
	// Bind correlation so downstream logs carry the expired pod's client identity.
	ctx = logger.WithCorrelation(ctx, logger.RequestIDFromContext(ctx), clientId, logger.TraceIDFromContext(ctx))

	logger.InfoEvent(ctx, "session_expired", "Expired VNC session detected",
		slog.String("session_id", pod.Name),
		slog.String("pod_name", pod.Name),
		slog.String("namespace", g.namespace),
	)

	metrics.RecordPodDeleteRequest()
	err := g.k8s.Clientset.CoreV1().Pods(g.namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
	if err != nil {
		_ = recordKubeFailure(ctx, "gc_delete_pod", lifecyclePodDelete, err,
			slog.String("session_id", pod.Name),
			slog.String("pod_name", pod.Name),
			slog.String("namespace", g.namespace),
		)
		return
	}

	metrics.RecordPodDeleteSuccess()
	logger.InfoEvent(ctx, "pod_delete_succeeded", "GC deleted expired VNC pod",
		slog.String("session_id", pod.Name),
		slog.String("pod_name", pod.Name),
		slog.String("namespace", g.namespace),
	)
}
