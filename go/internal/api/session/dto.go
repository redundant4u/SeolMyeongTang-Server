package session

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type getPodsResponse struct {
	Name      string `json:"name"`
	SessionId string `json:"sessionId"`
}

type createPodRequest struct {
	Name  string `json:"name" validate:"required,max=20"`
	Image string `json:"image" validate:"required,oneof=debian-xfce ubuntu"`
}

type createPodResponse struct {
	Name      string `json:"name"`
	SessionId string `json:"sessionId"`
}

type deletePodRequest struct {
	SessionId string `json:"sessionId" validate:"required"`
}

type createClientIdResponse struct {
	ClientId string `json:"clientId"`
}

func toGetSessionsResponse(pods []corev1.Pod) ([]getPodsResponse, error) {
	res := make([]getPodsResponse, 0, len(pods))

	for _, p := range pods {
		name, ok := p.Labels["name"]
		if !ok || name == "" {
			return nil, fmt.Errorf("pod label name is missing for pod: %s", p.Name)
		}

		res = append(res, getPodsResponse{
			Name:      name,
			SessionId: p.Name,
		})
	}

	return res, nil
}

func toCreateSessionResponse(pod *corev1.Pod, sessionId string) (createPodResponse, error) {
	name, ok := pod.Labels["name"]
	if !ok || name == "" {
		return createPodResponse{}, fmt.Errorf("pod label 'name' is missing")
	}

	return createPodResponse{
		Name:      name,
		SessionId: sessionId,
	}, nil
}
