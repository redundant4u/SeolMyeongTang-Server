package session

type getPodsResponse struct {
	Name      string `json:"name"`
	SessionId string `json:"sessionId"`
}

type createPodRequest struct {
	Name string
}

type createPodResponse struct {
	Name      string `json:"name"`
	SessionId string `json:"sessionId"`
}

type deletePodRequest struct {
	SessionId string `json:"sessionId"`
}

type createClientIdResponse struct {
	ClientId string `json:"clientId"`
}
