package metrics

import "github.com/prometheus/client_golang/prometheus"

var podCreateRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "smt_session_pod_create_requests_total",
		Help: "Total VNC pod creation requests.",
	},
	[]string{"image"},
)

var podCreateResults = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "smt_session_pod_create_results_total",
		Help: "VNC pod creation outcomes.",
	},
	[]string{"result", "error_code"},
)

var podDeleteRequests = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "smt_session_pod_delete_requests_total",
		Help: "Total VNC pod deletion requests.",
	},
)

var podDeleteResults = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "smt_session_pod_delete_results_total",
		Help: "VNC pod deletion outcomes.",
	},
	[]string{"result", "error_code"},
)

var activeSessions = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "smt_active_sessions",
		Help: "Current active VNC sessions observed by request handling.",
	},
)

func RecordPodCreateRequest(image string) {
	podCreateRequests.WithLabelValues(image).Inc()
}

func RecordPodCreateSuccess() {
	podCreateResults.WithLabelValues("success", "").Inc()
}

func RecordPodCreateFailure(errorCode string) {
	podCreateResults.WithLabelValues("failed", NormalizeErrorCode(errorCode)).Inc()
}

func RecordPodDeleteRequest() {
	podDeleteRequests.Inc()
}

func RecordPodDeleteSuccess() {
	podDeleteResults.WithLabelValues("success", "").Inc()
}

func RecordPodDeleteFailure(errorCode string) {
	podDeleteResults.WithLabelValues("failed", NormalizeErrorCode(errorCode)).Inc()
}

func SetActiveSessions(count int) {
	activeSessions.Set(float64(count))
}

func IncActiveSessions() {
	activeSessions.Inc()
}

func DecActiveSessions() {
	activeSessions.Dec()
}
