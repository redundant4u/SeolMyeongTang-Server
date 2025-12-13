package session

import (
	"context"
	"seolmyeong-tang-server/internal/pkg/k8s"
	"seolmyeong-tang-server/internal/pkg/logger"
	"time"

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

	logger.Info("gc: starting garbage collector")
	for {
		select {
		case <-ticker.C:
			g.cleanup(ctx)
		case <-ctx.Done():
			logger.Info("gc: stopping garbage collector")
			return
		}
	}
}

func (g *gc) cleanup(ctx context.Context) {
	listOptions := metav1.ListOptions{
		LabelSelector: "app=vnc",
		Limit:         100,
	}

	for {
		pods, err := g.k8s.Clientset.CoreV1().Pods(g.namespace).List(ctx, listOptions)
		if err != nil {
			logger.Error("gc: failed to list pods", err)
			return
		}

		now := time.Now().UTC()
		for _, pod := range pods.Items {
			if pod.ObjectMeta.DeletionTimestamp != nil {
				continue
			}

			expiredAtStr, ok := pod.Annotations["expired-at"]
			if !ok {
				continue
			}

			expiredAt, err := time.Parse(time.RFC3339, expiredAtStr)
			if err != nil {
				logger.Error("gc: failed to parse expired-at annotation", err)
				continue
			}

			if now.After(expiredAt) {
				err := g.k8s.Clientset.CoreV1().Pods(g.namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
				if err != nil {
					logger.Error("gc: failed to delete expired pod "+pod.Name, err)
				} else {
					logger.Info("gc: deleted expired pod " + pod.Name)
				}
			}
		}

		if pods.Continue == "" {
			break
		}

		listOptions.Continue = pods.Continue
	}
}
