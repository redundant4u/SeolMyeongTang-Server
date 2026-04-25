package session

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

const (
	cloudflaredImage     = "cloudflare/cloudflared:1800-17533b124c22"
	cloudflaredUserID    = int64(65532)
	runtimeClassName     = "kata-qemu-runtime-rs"
	sessionContainerPort = 5901
	sessionTTL           = 10 * time.Minute
)

// buildSessionPodSpec is a pure function that assembles the VNC session Pod
// spec. Extracted from kube.go so it can be unit-tested without a kube client.
func buildSessionPodSpec(info createPod, namespace string, now time.Time) *corev1.Pod {
	createdAt := now.UTC()
	expiredAt := createdAt.Add(sessionTTL)

	proxyURL := "http://vnc-gateway." + namespace + ".svc.cluster.local:3128"
	noProxyList := "localhost,127.0.0.1,.svc,.cluster.local,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

	sessionContainer := corev1.Container{
		Name:            info.sessionId,
		Image:           "vnc:" + info.image,
		ImagePullPolicy: corev1.PullNever,
		Env: []corev1.EnvVar{
			{Name: "HTTP_PROXY", Value: proxyURL},
			{Name: "http_proxy", Value: proxyURL},
			{Name: "HTTPS_PROXY", Value: proxyURL},
			{Name: "https_proxy", Value: proxyURL},
			{Name: "NO_PROXY", Value: noProxyList},
			{Name: "no_proxy", Value: noProxyList},
		},
		Ports: []corev1.ContainerPort{
			{Name: "vnc", ContainerPort: sessionContainerPort},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: "workspace", MountPath: "/home/app"},
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:              resource.MustParse("500m"),
				corev1.ResourceMemory:           resource.MustParse("1Gi"),
				corev1.ResourceEphemeralStorage: resource.MustParse("3Gi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:              resource.MustParse("1"),
				corev1.ResourceMemory:           resource.MustParse("2Gi"),
				corev1.ResourceEphemeralStorage: resource.MustParse("5Gi"),
			},
		},
	}

	cloudflaredContainer := corev1.Container{
		Name:  "cloudflared",
		Image: cloudflaredImage,
		Args: []string{
			"tunnel",
			"--config",
			"/etc/cloudflared/config.yml",
			"--no-autoupdate",
			"run",
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "cloudflared-config",
				MountPath: "/etc/cloudflared/config.yml",
				SubPath:   "config.yml",
			},
			{
				Name:      "cloudflared-credentials",
				MountPath: "/etc/cloudflared/credentials.json",
				SubPath:   "credentials.json",
			},
		},
		SecurityContext: &corev1.SecurityContext{
			RunAsNonRoot:             ptr.To(true),
			RunAsUser:                ptr.To(cloudflaredUserID),
			AllowPrivilegeEscalation: ptr.To(false),
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("50m"),
				corev1.ResourceMemory: resource.MustParse("64Mi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("200m"),
				corev1.ResourceMemory: resource.MustParse("256Mi"),
			},
		},
	}

	workspaceSize := resource.MustParse("5Gi")

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      info.sessionId,
			Labels: map[string]string{
				"app":       "vnc",
				"name":      info.name,
				"client-id": info.clientId,
			},
			Annotations: map[string]string{
				"description": info.description,
				"created-at":  createdAt.Format(time.RFC3339),
				"expired-at":  expiredAt.Format(time.RFC3339),
			},
		},
		Spec: corev1.PodSpec{
			AutomountServiceAccountToken: ptr.To(false),
			RestartPolicy:                corev1.RestartPolicyNever,
			RuntimeClassName:             ptr.To(runtimeClassName),
			Volumes: []corev1.Volume{
				{
					Name: "workspace",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							SizeLimit: &workspaceSize,
						},
					},
				},
				{
					Name: "cloudflared-credentials",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "cloudflared-credentials",
						},
					},
				},
				{
					Name: "cloudflared-config",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "cloudflared-config-" + info.sessionId,
							},
						},
					},
				},
			},
			Containers: []corev1.Container{sessionContainer, cloudflaredContainer},
		},
	}
}
