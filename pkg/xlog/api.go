// Package xlog logging module
package xlog

import (
	"context"

	"go.uber.org/zap"
)

// Info logs a message at Info level with context fields.
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	log.Info(msg, append(fields, ctxFields(ctx)...)...)
}

// Debug logs a message at Debug level with context fields.
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	log.Debug(msg, append(fields, ctxFields(ctx)...)...)
}

// Warn logs a message at Warn level with context fields.
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	log.Warn(msg, append(fields, ctxFields(ctx)...)...)
}

// Error logs a message at Error level with context fields.
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	log.Error(msg, append(fields, ctxFields(ctx)...)...)
}

// Fatal logs a message at Fatal level with context fields, then calls os.Exit(1).
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	log.Fatal(msg, append(fields, ctxFields(ctx)...)...)
}

// InfoNoCtx logs a message at Info level without context fields.
func InfoNoCtx(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// ErrorNoCtx logs a message at Error level without context fields.
func ErrorNoCtx(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// FatalNoCtx logs a message at Fatal level without context fields, then calls os.Exit(1).
func FatalNoCtx(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}
