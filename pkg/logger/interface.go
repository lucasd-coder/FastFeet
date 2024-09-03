package logger

import (
	"context"
	"log/slog"
)

type Fields map[string]interface{}

type Logger interface {
	Info(msg string, args ...any)

	Debug(msg string, args ...any)

	Error(msg string, args ...any)

	Warn(msg string, args ...any)

	Trace(msg string, args ...any)

	Fatal(msg string, args ...any)

	Debugf(format string, args ...any)

	Infof(format string, args ...any)

	Warnf(format string, args ...any)

	Errorf(format string, args ...any)

	Tracef(format string, args ...any)

	Fatalf(format string, args ...any)

	With(args ...any) Logger

	LogLevel(ctx context.Context, level slog.Level, msg string, args ...any)

	LogLevelf(ctx context.Context, level slog.Level, format string, args ...any)
}
