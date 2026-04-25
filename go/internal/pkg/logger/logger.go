package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

const ServiceName = "smt-server"

var base = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
		switch attr.Key {
		case slog.TimeKey:
			attr.Key = "timestamp"
		case slog.MessageKey:
			attr.Key = "message"
		case slog.LevelKey:
			level := attr.Value.String()
			attr.Value = slog.StringValue(strings.ToLower(level))
		}
		return attr
	},
})).With(
	slog.String("service", ServiceName),
)

func Event(ctx context.Context, level slog.Level, eventType, message string, attrs ...slog.Attr) {
	values := make([]any, 0, len(attrs)+5)
	values = append(values, slog.String("event_type", eventType))
	for _, attr := range CorrelationAttrs(ctx) {
		values = append(values, attr)
	}
	for _, attr := range SafeAttrs(attrs...) {
		values = append(values, attr)
	}
	base.Log(ctx, level, message, values...)
}

func InfoEvent(ctx context.Context, eventType, message string, attrs ...slog.Attr) {
	Event(ctx, slog.LevelInfo, eventType, message, attrs...)
}

func WarnEvent(ctx context.Context, eventType, message string, attrs ...slog.Attr) {
	Event(ctx, slog.LevelWarn, eventType, message, attrs...)
}

func ErrorEvent(ctx context.Context, eventType, message string, err error, attrs ...slog.Attr) {
	if err != nil {
		attrs = append(attrs, slog.String("error", RedactString(err.Error())))
	}
	Event(ctx, slog.LevelError, eventType, message, attrs...)
}

func FatalEvent(ctx context.Context, eventType, message string, err error, attrs ...slog.Attr) {
	ErrorEvent(ctx, eventType, message, err, attrs...)
	os.Exit(1)
}
