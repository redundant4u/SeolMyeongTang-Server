package logger

import (
	"context"
	"log/slog"
	"testing"
)

func TestSafeAttrsRedactsSensitiveFields(t *testing.T) {
	attrs := SafeAttrs(
		slog.String("password", "plain"),
		slog.String("authorization", "Bearer abc"),
		slog.String("message", "ok"),
	)

	if got := attrs[0].Value.String(); got != "[REDACTED]" {
		t.Fatalf("password was not redacted: %q", got)
	}
	if got := attrs[1].Value.String(); got != "[REDACTED]" {
		t.Fatalf("authorization was not redacted: %q", got)
	}
	if got := attrs[2].Value.String(); got != "ok" {
		t.Fatalf("non-sensitive value changed: %q", got)
	}
}

func TestClientDisplayIDShortensLongIdentifiers(t *testing.T) {
	got := ClientDisplayID("abcdefgh12345678")
	if got != "abcd...5678" {
		t.Fatalf("unexpected display id: %q", got)
	}
}

func TestRedactStringScrubsNovncAndBearerMaterial(t *testing.T) {
	cases := []string{
		"failed to reach https://redundant4u.com/?novnc_token=abc123",
		"Authorization: Bearer eyJhbGciOi",
		"api_key=supersecret",
		"bearer xyz",
	}
	for _, in := range cases {
		if got := RedactString(in); got == in {
			t.Fatalf("expected redaction for %q, got %q", in, got)
		}
	}
}

func TestSafeAttrsPreservesCorrelationAttrs(t *testing.T) {
	attrs := SafeAttrs(
		slog.String("request_id", "req-abcdefabcdefabcdefabcdefabcdef12"),
		slog.String("client_id", "client-123"),
		slog.String("session_id", "sess-xyz"),
	)
	if got := attrs[0].Value.String(); got != "req-abcdefabcdefabcdefabcdefabcdef12" {
		t.Fatalf("request_id should not be redacted: %q", got)
	}
	if got := attrs[1].Value.String(); got != "client-123" {
		t.Fatalf("client_id should not be redacted: %q", got)
	}
	if got := attrs[2].Value.String(); got != "sess-xyz" {
		t.Fatalf("session_id should not be redacted: %q", got)
	}
}

func TestCorrelationAttrsFromContext(t *testing.T) {
	ctx := WithCorrelation(context.Background(), "req-1", "client-abc", "trace-x")
	attrs := CorrelationAttrs(ctx)

	attrMap := make(map[string]string, len(attrs))
	for _, a := range attrs {
		attrMap[a.Key] = a.Value.String()
	}

	if attrMap["request_id"] != "req-1" {
		t.Fatalf("request_id mismatch: %q", attrMap["request_id"])
	}
	if attrMap["client_id"] != "client-abc" {
		t.Fatalf("client_id mismatch: %q", attrMap["client_id"])
	}
	if attrMap["client_display_id"] != "clie...-abc" {
		t.Fatalf("client_display_id mismatch: %q", attrMap["client_display_id"])
	}
	if attrMap["trace_id"] != "trace-x" {
		t.Fatalf("trace_id mismatch: %q", attrMap["trace_id"])
	}
}

func TestCorrelationAttrsOmitsEmptyValues(t *testing.T) {
	ctx := WithCorrelation(context.Background(), "", "", "")
	attrs := CorrelationAttrs(ctx)
	if len(attrs) != 0 {
		t.Fatalf("expected no attrs for empty correlation, got %d", len(attrs))
	}
}
