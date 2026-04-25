package metrics

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var registerOnce sync.Once

// HTTP request counter carries an outcome label; the histogram intentionally
// omits outcome to keep latency series cardinality bounded.
var httpRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "smt_http_requests_total",
		Help: "Total HTTP requests handled by the SMT server.",
	},
	[]string{"method", "path", "status", "outcome"},
)

var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "smt_http_request_duration_seconds",
		Help:    "HTTP request latency for the SMT server.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path", "status"},
)

var kubeAPIFailures = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "smt_kubernetes_api_failures_total",
		Help: "Kubernetes API failures observed by the SMT server.",
	},
	[]string{"operation", "error_code"},
)

func Register() {
	registerOnce.Do(func() {
		prometheus.MustRegister(
			httpRequests,
			httpDuration,
			kubeAPIFailures,
			podCreateRequests,
			podCreateResults,
			podDeleteRequests,
			podDeleteResults,
			activeSessions,
		)
	})
}

func Handler() http.Handler {
	Register()
	return promhttp.Handler()
}

func RecordHTTPRequest(method, path, status, outcome string, latency time.Duration) {
	httpRequests.WithLabelValues(method, path, status, outcome).Inc()
	httpDuration.WithLabelValues(method, path, status).Observe(latency.Seconds())
}

func RecordKubeAPIFailure(operation, errorCode string) {
	kubeAPIFailures.WithLabelValues(operation, NormalizeErrorCode(errorCode)).Inc()
}

func StatusLabel(status int) string {
	return strconv.Itoa(status)
}

// NormalizeErrorCode replaces empty codes with "unknown" so no metric series
// carry an empty error_code label.
func NormalizeErrorCode(code string) string {
	if code == "" {
		return "unknown"
	}
	return code
}
