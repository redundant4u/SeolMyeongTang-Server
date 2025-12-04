package session

type getPodsResponse struct {
	Name      string `json:"name"`
	SessionId string `json:"sessionId"`
}

type createPodRequest struct {
	Name      string
	SessionId string
}

type createPodResponse struct {
	Name      string `json:"name"`
	SessionId string `json:"sessionId"`
}

type deletePodRequest struct {
	SessionId string `json:"sessionId"`
}
