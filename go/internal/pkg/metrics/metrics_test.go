package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestRecordHTTPRequestExposesMetric(t *testing.T) {
	counter := httpRequests.WithLabelValues("GET", "/metric-test", "200", "success")
	histogram := httpDuration.WithLabelValues("GET", "/metric-test", "200").(prometheus.Metric)
	beforeCount := readMetric(t, counter).GetCounter().GetValue()
	beforeLatency := readMetric(t, histogram).GetHistogram()
	RecordHTTPRequest("GET", "/metric-test", "200", "success", 10*time.Millisecond)
	if got := readMetric(t, counter).GetCounter().GetValue(); got != beforeCount+1 {
		t.Fatalf("request counter did not increment: %f", got)
	}
	after := readMetric(t, histogram).GetHistogram()
	if after.GetSampleCount() != beforeLatency.GetSampleCount()+1 || after.GetSampleSum() != beforeLatency.GetSampleSum()+0.01 {
		t.Fatal("request latency was not observed")
	}

	if !metricFamilyExists("smt_http_requests_total") {
		t.Fatal("smt_http_requests_total was not gathered")
	}
	if !metricFamilyExists("smt_http_request_duration_seconds") {
		t.Fatal("smt_http_request_duration_seconds was not gathered")
	}
}

func TestSessionMetricHelpersExposeLifecycleMetrics(t *testing.T) {
	for _, tc := range []struct {
		counter prometheus.Counter
		record  func()
	}{
		{podCreateRequests.WithLabelValues("ubuntu"), func() { RecordPodCreateRequest("ubuntu") }},
		{podCreateResults.WithLabelValues("success", ""), RecordPodCreateSuccess},
		{podCreateResults.WithLabelValues("failed", "unknown"), func() { RecordPodCreateFailure("") }},
		{podDeleteRequests, RecordPodDeleteRequest},
		{podDeleteResults.WithLabelValues("success", ""), RecordPodDeleteSuccess},
		{podDeleteResults.WithLabelValues("failed", "denied"), func() { RecordPodDeleteFailure("denied") }},
	} {
		before := readMetric(t, tc.counter).GetCounter().GetValue()
		tc.record()
		if got := readMetric(t, tc.counter).GetCounter().GetValue(); got != before+1 {
			t.Fatal("lifecycle counter did not increment")
		}
	}
	SetActiveSessions(2)
	IncActiveSessions()
	if readMetric(t, activeSessions).GetGauge().GetValue() != 3 {
		t.Fatal("active session increment failed")
	}
	DecActiveSessions()
	if readMetric(t, activeSessions).GetGauge().GetValue() != 2 {
		t.Fatal("active session decrement failed")
	}

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

func readMetric(t *testing.T, metric prometheus.Metric) *dto.Metric {
	t.Helper()
	value := &dto.Metric{}
	if err := metric.Write(value); err != nil {
		t.Fatal(err)
	}
	return value
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
