package logger

import (
	"log/slog"
	"regexp"
	"strings"
)

var sensitiveFragments = []string{
	"password",
	"passwd",
	"token",
	"cookie",
	"secret",
	"credential",
	"authorization",
	"bearer",
	"api_key",
	"apikey",
	"session_value",
	"novnc",
	"vnc_token",
}

// Matches `<fragment>[=: ]<value>` forms such as `token=abc`, `Authorization: Bearer ...`,
// or `bearer abc`. Used to scrub free-form error messages from kube/network libraries
// that may echo secret material.
var inlineSensitivePattern = regexp.MustCompile(
	`(?i)(` + strings.Join(sensitiveFragments, "|") + `)\s*[:=]?\s*[^\s,;]+`,
)

// Attr keys used for correlation identifiers. Their values must never be pattern-redacted
// because dashboards and log correlation depend on them staying intact.
var correlationAttrKeys = map[string]struct{}{
	"request_id":        {},
	"trace_id":          {},
	"client_id":         {},
	"client_display_id": {},
	"session_id":        {},
	"pod_name":          {},
	"namespace":         {},
}

func SafeAttrs(attrs ...slog.Attr) []slog.Attr {
	safe := make([]slog.Attr, 0, len(attrs))
	for _, attr := range attrs {
		key := strings.ToLower(attr.Key)
		if isSensitiveKey(key) {
			safe = append(safe, slog.String(attr.Key, "[REDACTED]"))
			continue
		}
		if _, ok := correlationAttrKeys[key]; ok {
			safe = append(safe, attr)
			continue
		}
		if attr.Value.Kind() == slog.KindString {
			safe = append(safe, slog.String(attr.Key, RedactString(attr.Value.String())))
			continue
		}
		safe = append(safe, attr)
	}
	return safe
}

// RedactString scrubs free-form text (typically error messages) for sensitive material.
// Structured attr values for known correlation keys bypass this path via SafeAttrs.
func RedactString(value string) string {
	return inlineSensitivePattern.ReplaceAllString(value, "[REDACTED]")
}

func ClientDisplayID(clientID string) string {
	if len(clientID) <= 8 {
		return clientID
	}
	return clientID[:4] + "..." + clientID[len(clientID)-4:]
}

func isSensitiveKey(key string) bool {
	for _, fragment := range sensitiveFragments {
		if strings.Contains(key, fragment) {
			return true
		}
	}
	return false
}
