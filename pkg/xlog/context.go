package xlog

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

// WithRequestID stores a request ID in the context for structured logging.
func WithRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, requestIDKey, rid)
}

// RequestID retrieves the request ID from the context, returning an empty string if absent.
func RequestID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func ctxFields(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	if rid := RequestID(ctx); rid != "" {
		return []zap.Field{zap.String("rid", rid)}
	}

	return nil
}
