package session

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"seolmyeong-tang-server/internal/config"
	"seolmyeong-tang-server/internal/pkg/k8s"
	"seolmyeong-tang-server/internal/pkg/logger"
	"seolmyeong-tang-server/internal/pkg/metrics"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type Kube struct {
	k8s       *k8s.Client
	namespace string
	Gc        *gc
}

type createPod struct {
	name        string
	image       string
	clientId    string
	sessionId   string
	description string
}

type deletePod struct {
	clientId  string
	sessionId string
}

func NewKube(k8s *k8s.Client, namespace string) *Kube {
	return &Kube{
		k8s:       k8s,
		namespace: namespace,
		Gc:        newGC(k8s, namespace, 60),
	}
}

func (k *Kube) getPods(ctx context.Context, clientId string) ([]corev1.Pod, error) {
	selector := labels.Set{
		"app":       "vnc",
		"client-id": clientId,
	}.AsSelector().String()

	pods, err := k.k8s.Clientset.CoreV1().
		Pods(k.namespace).
		List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, recordKubeFailure(ctx, "list_pods", lifecycleNone, err,
			slog.String("namespace", k.namespace),
		)
	}

	running := make([]corev1.Pod, 0, len(pods.Items))
	for _, p := range pods.Items {
		if p.ObjectMeta.DeletionTimestamp != nil {
			continue
		}
		if p.Status.Phase == corev1.PodRunning || p.Status.Phase == corev1.PodPending {
			running = append(running, p)
		}
	}

	return running, nil
}

func (k *Kube) getSessions(ctx context.Context, clientId string) ([]corev1.Pod, error) {
	return k.getPods(ctx, clientId)
}

func (k *Kube) createSession(ctx context.Context, info createPod) (*corev1.Pod, error) {
	logger.InfoEvent(ctx, "pod_create_started", "VNC pod creation started",
		slog.String("session_id", info.sessionId),
		slog.String("pod_name", info.sessionId),
		slog.String("namespace", k.namespace),
		slog.String("requested_image", info.image),
	)

	if err := k.createCloudflaredConfigMap(ctx, info.sessionId); err != nil {
		return nil, recordKubeFailure(ctx, "create_configmap", lifecyclePodCreate, err,
			slog.String("session_id", info.sessionId),
			slog.String("pod_name", info.sessionId),
			slog.String("namespace", k.namespace),
		)
	}

	pod, err := k.k8s.Clientset.CoreV1().
		Pods(k.namespace).
		Create(ctx, buildSessionPodSpec(info, k.namespace, time.Now()), metav1.CreateOptions{})
	if err != nil {
		return nil, recordKubeFailure(ctx, "create_pod", lifecyclePodCreate, err,
			slog.String("session_id", info.sessionId),
			slog.String("pod_name", info.sessionId),
			slog.String("namespace", k.namespace),
		)
	}

	metrics.RecordPodCreateSuccess()
	metrics.IncActiveSessions()
	logger.InfoEvent(ctx, "pod_create_succeeded", "VNC pod creation succeeded",
		slog.String("session_id", info.sessionId),
		slog.String("pod_name", pod.Name),
		slog.String("namespace", k.namespace),
	)

	return pod, nil
}

func (k *Kube) deleteSession(ctx context.Context, info deletePod) error {
	if err := k.deleteCloudflaredConfigMap(ctx, info.sessionId); err != nil {
		return recordKubeFailure(ctx, "delete_configmap", lifecyclePodDelete, err,
			slog.String("session_id", info.sessionId),
			slog.String("namespace", k.namespace),
		)
	}

	pods, err := k.getPods(ctx, info.clientId)
	if err != nil {
		metrics.RecordPodDeleteFailure(kubeErrorCode(err))
		return err
	}

	var target *corev1.Pod
	for i := range pods {
		if pods[i].Name == info.sessionId {
			target = &pods[i]
			break
		}
	}

	if target == nil {
		logger.WarnEvent(ctx, "pod_delete_skipped", "VNC pod delete skipped because pod was not found",
			slog.String("session_id", info.sessionId),
			slog.String("pod_name", info.sessionId),
			slog.String("namespace", k.namespace),
		)
		return nil
	}

	logger.InfoEvent(ctx, "pod_delete_started", "VNC pod deletion started",
		slog.String("session_id", info.sessionId),
		slog.String("pod_name", info.sessionId),
		slog.String("namespace", k.namespace),
	)

	if err := k.k8s.Clientset.CoreV1().
		Pods(k.namespace).
		Delete(ctx, info.sessionId, metav1.DeleteOptions{}); err != nil {
		return recordKubeFailure(ctx, "delete_pod", lifecyclePodDelete, err,
			slog.String("session_id", info.sessionId),
			slog.String("pod_name", info.sessionId),
			slog.String("namespace", k.namespace),
		)
	}

	metrics.RecordPodDeleteSuccess()
	metrics.DecActiveSessions()
	logger.InfoEvent(ctx, "pod_delete_succeeded", "VNC pod deletion succeeded",
		slog.String("session_id", info.sessionId),
		slog.String("pod_name", info.sessionId),
		slog.String("namespace", k.namespace),
	)

	return nil
}

func (k *Kube) createCloudflaredConfigMap(ctx context.Context, sessionId string) error {
	configMapName := "cloudflared-config-" + sessionId
	// config.yml는 공백 기반 포맷이라 탭 문자를 넣으면 안 됨.
	config := fmt.Sprintf(`
tunnel: %s
credentials-file: /etc/cloudflared/credentials.json

ingress:
  - hostname: %s.tunnel.redundant4u.com
    service: http://localhost:8080
  - service: http_status:404`, config.Env.CF_TUNNEL_ID, sessionId)

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: k.namespace,
			Name:      configMapName,
		},
		Data: map[string]string{"config.yml": config},
	}

	_, err := k.k8s.Clientset.CoreV1().
		ConfigMaps(k.namespace).
		Create(ctx, configMap, metav1.CreateOptions{})
	return unwrapKubeStatus(err)
}

func (k *Kube) deleteCloudflaredConfigMap(ctx context.Context, sessionId string) error {
	err := k.k8s.Clientset.CoreV1().
		ConfigMaps(k.namespace).
		Delete(ctx, "cloudflared-config-"+sessionId, metav1.DeleteOptions{})
	return unwrapKubeStatus(err)
}

func (k *Kube) secureRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}
