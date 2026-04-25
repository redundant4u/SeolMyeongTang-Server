package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestRecordHTTPRequestExposesMetric(t *testing.T) {
	RecordHTTPRequest("GET", "/session", "200", "success", 10*time.Millisecond)

	if !metricFamilyExists("smt_http_requests_total") {
		t.Fatal("smt_http_requests_total was not gathered")
	}
	if !metricFamilyExists("smt_http_request_duration_seconds") {
		t.Fatal("smt_http_request_duration_seconds was not gathered")
	}
}

func TestSessionMetricHelpersExposeLifecycleMetrics(t *testing.T) {
	RecordPodCreateRequest("ubuntu")
	RecordPodCreateSuccess()
	RecordPodDeleteRequest()
	RecordPodDeleteSuccess()
	SetActiveSessions(2)

	for _, name := range []string{
		"smt_session_pod_create_requests_total",
		"smt_session_pod_create_results_total",
		"smt_session_pod_delete_requests_total",
		"smt_session_pod_delete_results_total",
		"smt_active_sessions",
	} {
		if !metricFamilyExists(name) {
			t.Fatalf("%s was not gathered", name)
		}
	}
}

func metricFamilyExists(name string) bool {
	Register()
	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return false
	}
	for _, family := range families {
		if family.GetName() == name {
			return true
		}
	}
	return false
}
