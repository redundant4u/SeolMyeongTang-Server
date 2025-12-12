package session

import (
	"context"
	"crypto/rand"
	"fmt"
	"seolmyeong-tang-server/internal/pkg/k8s"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type kube struct {
	k8s       *k8s.Client
	namespace string
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

func newKube(k8s *k8s.Client, namespace string) *kube {
	return &kube{k8s: k8s, namespace: namespace}
}

func (k *kube) getPods(ctx context.Context, clientId string) ([]corev1.Pod, error) {
	selector := labels.Set{
		"app":       "vnc",
		"client-id": clientId,
	}.AsSelector().String()

	pods, err := k.k8s.Clientset.CoreV1().
		Pods("vnc").
		List(ctx, metav1.ListOptions{
			LabelSelector: selector,
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

func (k *kube) getSessions(ctx context.Context, clientId string) ([]corev1.Pod, error) {
	return k.getPods(ctx, clientId)
}

func (k *kube) createSession(ctx context.Context, info createPod) (*corev1.Pod, error) {
	podSpec := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "vnc",
			Name:      info.sessionId,
			Labels: map[string]string{
				"app":       "vnc",
				"name":      info.name,
				"client-id": info.clientId,
			},
			Annotations: map[string]string{
				"description": info.description,
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Volumes: []corev1.Volume{
				{
					Name: "workspace",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							SizeLimit: func() *resource.Quantity {
								q := resource.MustParse("5Gi")
								return &q
							}(),
						},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:            info.sessionId,
					Image:           "vnc:" + info.image,
					ImagePullPolicy: "Never",
					Ports: []corev1.ContainerPort{
						{
							Name:          "vnc",
							ContainerPort: 5901,
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "workspace",
							MountPath: "/home/app",
						},
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:              resource.MustParse("500m"),
							corev1.ResourceMemory:           resource.MustParse("512Mi"),
							corev1.ResourceEphemeralStorage: resource.MustParse("3Gi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:              resource.MustParse("1"),
							corev1.ResourceMemory:           resource.MustParse("1Gi"),
							corev1.ResourceEphemeralStorage: resource.MustParse("5Gi"),
						},
					},
				},
			},
		},
	}

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

func (k *kube) deleteSession(ctx context.Context, info deletePod) error {
	pods, err := k.getPods(ctx, info.clientId)
	if err != nil {
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
		return nil
	}

	if err := k.k8s.Clientset.CoreV1().
		Pods(k.namespace).
		Delete(ctx, info.sessionId, metav1.DeleteOptions{}); err != nil {

		if statusErr, ok := err.(kerrors.APIStatus); ok {
			st := statusErr.Status()
			return fmt.Errorf("code=%d reason=%s message=%s", st.Code, st.Reason, st.Message)
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
