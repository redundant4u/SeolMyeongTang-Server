package logger

import (
	"context"
	"log/slog"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	clientIDKey  contextKey = "client_id"
	traceIDKey   contextKey = "trace_id"
)

func WithCorrelation(ctx context.Context, requestID, clientID, traceID string) context.Context {
	ctx = context.WithValue(ctx, requestIDKey, requestID)
	ctx = context.WithValue(ctx, clientIDKey, clientID)
	ctx = context.WithValue(ctx, traceIDKey, traceID)
	return ctx
}

func RequestIDFromContext(ctx context.Context) string {
	return stringValue(ctx, requestIDKey)
}

func ClientIDFromContext(ctx context.Context) string {
	return stringValue(ctx, clientIDKey)
}

func TraceIDFromContext(ctx context.Context) string {
	return stringValue(ctx, traceIDKey)
}

// CorrelationAttrs returns non-empty request_id/client_id/client_display_id/trace_id
// attributes pulled from ctx. Empty identifiers are omitted so logs don't carry
// pointless empty strings.
func CorrelationAttrs(ctx context.Context) []slog.Attr {
	if ctx == nil {
		return nil
	}
	attrs := make([]slog.Attr, 0, 4)
	if v := stringValue(ctx, requestIDKey); v != "" {
		attrs = append(attrs, slog.String("request_id", v))
	}
	if v := stringValue(ctx, clientIDKey); v != "" {
		attrs = append(attrs, slog.String("client_id", v))
		attrs = append(attrs, slog.String("client_display_id", ClientDisplayID(v)))
	}
	if v := stringValue(ctx, traceIDKey); v != "" {
		attrs = append(attrs, slog.String("trace_id", v))
	}
	return attrs
}

func stringValue(ctx context.Context, key contextKey) string {
	value, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return value
}
