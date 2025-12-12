package session

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type getPodsResponse struct {
	Name        string `json:"name"`
	SessionId   string `json:"sessionId"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type createPodRequest struct {
	Name        string `json:"name" validate:"required,k8slabel,max=20"`
	Image       string `json:"image" validate:"required,oneof=debian-xfce ubuntu"`
	Description string `json:"description" validate:"k8slabel,max=100"`
}

type createPodResponse struct {
	Name        string `json:"name"`
	SessionId   string `json:"sessionId"`
	Image       string `json:"image"`
	Description string `json:"description"`
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
			return nil, fmt.Errorf("pod label 'name' is missing")
		}

		image := p.Spec.Containers[0].Image

		description, ok := p.Annotations["description"]
		if !ok {
			description = ""
		}

		res = append(res, getPodsResponse{
			Name:        name,
			SessionId:   p.Name,
			Image:       image,
			Description: description,
		})
	}

	return res, nil
}

func toCreateSessionResponse(pod *corev1.Pod, sessionId string) (createPodResponse, error) {
	name, ok := pod.Labels["name"]
	if !ok || name == "" {
		return createPodResponse{}, fmt.Errorf("pod label 'name' is missing")
	}

	image := pod.Spec.Containers[0].Image

	description, ok := pod.Annotations["description"]
	if !ok {
		description = ""
	}

	return createPodResponse{
		Name:        name,
		SessionId:   sessionId,
		Image:       image,
		Description: description,
	}, nil
}
