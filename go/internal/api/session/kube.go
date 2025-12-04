package session

import (
	"context"
	"crypto/rand"
	"fmt"
	"seolmyeong-tang-server/internal/pkg/k8s"
	"seolmyeong-tang-server/internal/pkg/logger"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type kube struct {
	k8s       *k8s.Client
	namespace string
}

func newKube(k8s *k8s.Client, namespace string) *kube {
	return &kube{k8s: k8s, namespace: namespace}
}

func (k *kube) getSessions(ctx context.Context) ([]corev1.Pod, error) {
	pods, err := k.k8s.Clientset.CoreV1().
		Pods("vnc").
		List(ctx, metav1.ListOptions{
			LabelSelector: "app=vnc",
		})
	if err != nil {
		return nil, err
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

	return running, err
}

func (k *kube) createSession(ctx context.Context, info createPodRequest) (*corev1.Pod, error) {
	podSpec := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "vnc",
			Name:      info.Name + "-" + info.SessionId,
			Labels: map[string]string{
				"app":        "vnc",
				"session-id": info.SessionId,
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:            "vnc-server",
					Image:           "vnc",
					ImagePullPolicy: "Never",
					Ports: []corev1.ContainerPort{
						{
							Name:          "vnc",
							ContainerPort: 5901,
						},
					},
				},
			},
		},
	}

	logger.Info(k.namespace)
	pod, err := k.k8s.Clientset.CoreV1().
		Pods(k.namespace).
		Create(ctx, podSpec, metav1.CreateOptions{})
	if err != nil {
		if statusErr, ok := err.(kerrors.APIStatus); ok {
			st := statusErr.Status()

			return nil, fmt.Errorf(
				"code=%d reason=%s message=%s",
				st.Code,
				st.Reason,
				st.Message,
			)
		}

		return nil, err
	}

	return pod, nil
}

func (k *kube) deleteSession(ctx context.Context, info deletePodRequest) error {
	podName := fmt.Sprintf("vnc-%s", info.SessionId)

	err := k.k8s.Clientset.CoreV1().
		Pods(k.namespace).
		Delete(ctx, podName, metav1.DeleteOptions{})

	if err != nil {
		if statusErr, ok := err.(kerrors.APIStatus); ok {
			st := statusErr.Status()

			return fmt.Errorf(
				"code=%d reason=%s message=%s",
				st.Code,
				st.Reason,
				st.Message,
			)
		}

		return err
	}
	return nil
}

func (k *kube) secureRandomString(n int) string {
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
