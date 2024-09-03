package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const (
	LevelTrace   = slog.Level(-8)
	LevelDebug   = slog.LevelDebug
	LevelInfo    = slog.LevelInfo
	LevelNotice  = slog.Level(2)
	LevelWarning = slog.LevelWarn
	LevelError   = slog.LevelError
	LevelFatal   = slog.Level(12)
)

var opts = []logging.Option{
	logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
	logging.WithFieldsFromContext(logTraceID),
}

var logTraceID = func(ctx context.Context) logging.Fields {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return logging.Fields{"traceID", span.TraceID().String()}
	}
	return nil
}

type Option struct {
	AppName      string
	Level        string
	ReportCaller bool
}

type loggingConfig struct {
	mu          sync.RWMutex
	l           *Log
	slogDefault *slog.Logger
}

var config = &loggingConfig{}

type Log struct {
	opt Option
}

func NewLogger(opt Option) *Log {
	config.mu.Lock()
	defer config.mu.Unlock()
	if config.l == nil {
		config.l = &Log{opt: opt}
	}
	return config.l
}

func (l *Log) GetLog() *slog.Logger {
	config.mu.Lock()
	defer config.mu.Unlock()

	if config.slogDefault == nil {
		config.slogDefault = l.createLogger()
	}

	return config.slogDefault
}

func (l *Log) createLogger() *slog.Logger {
	level := l.parseLogLevel()
	opts := &slog.HandlerOptions{
		Level:       level,
		AddSource:   l.opt.ReportCaller,
		ReplaceAttr: l.replaceAttr,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler).
		With(slog.String("application", l.opt.AppName))
}

func (l *Log) replaceAttr(groups []string, a slog.Attr) slog.Attr {
	// Remove time from the output for predictable test output.
	if a.Key == slog.TimeKey && len(groups) == 0 {
		formattedTime := a.Value.Any().(time.Time).Format(time.RFC3339)
		return slog.Attr{
			Key:   slog.TimeKey,
			Value: slog.StringValue(formattedTime),
		}
	}
	// Customize the name of the level key and the output string, including
	// custom level values.
	if a.Key == slog.LevelKey {
		// Rename the level key from "level" to "sev"
		// Handle custom level values.
		level := a.Value.Any().(slog.Level)

		// This could also look up the name from a map or other structure, but
		// this demonstrates using a switch statement to rename levels. For
		// maximum performance, the string values should be constants, but this
		// example uses the raw strings for readability.
		switch {
		case level < LevelDebug:
			a.Value = slog.StringValue("TRACE")
		case level < LevelInfo:
			a.Value = slog.StringValue("DEBUG")
		case level < LevelNotice:
			a.Value = slog.StringValue("INFO")
		case level < LevelWarning:
			a.Value = slog.StringValue("NOTICE")
		case level < LevelError:
			a.Value = slog.StringValue("WARNING")
		case level < LevelFatal:
			a.Value = slog.StringValue("ERROR")
		default:
			a.Value = slog.StringValue("FATAL")
		}
	}
	return a
}

func (l *Log) parseLogLevel() slog.Level {
	switch l.opt.Level {
	case "INFO":
		return LevelInfo
	case "ERROR":
		return LevelError
	case "DEBUG":
		return LevelDebug
	case "WARN":
		return LevelWarning
	case "NOTICE":
		return LevelNotice
	case "FATAL":
		return LevelFatal
	default:
		return LevelTrace
	}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (l *Log) GetLogUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(InterceptorLogger(slog.Default()), opts...)
}

func (l *Log) GetLogStreamServerInterceptor() grpc.StreamServerInterceptor {
	return logging.StreamServerInterceptor(InterceptorLogger(slog.Default()), opts...)
}

func (l *Log) GetLogUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return logging.UnaryClientInterceptor(InterceptorLogger(slog.Default()), opts...)
}

func (l *Log) GetLogStreamClientInterceptor() grpc.StreamClientInterceptor {
	return logging.StreamClientInterceptor(InterceptorLogger(slog.Default()), opts...)
}

type log struct {
	*slog.Logger
	ctx context.Context
}

func (l log) Info(msg string, args ...any) {
	l.InfoContext(l.ctx, msg, args...)
}

func (l log) Infof(format string, args ...any) {
	l.InfoContext(l.ctx, fmt.Sprintf(format, args...))
}

func (l log) Warn(msg string, args ...any) {
	l.WarnContext(l.ctx, msg, args...)
}

func (l log) Warnf(format string, args ...any) {
	l.WarnContext(l.ctx, fmt.Sprintf(format, args...))
}

func (l log) Debug(msg string, args ...any) {
	l.DebugContext(l.ctx, msg, args...)
}

func (l log) Debugf(format string, args ...any) {
	l.DebugContext(l.ctx, fmt.Sprintf(format, args...))
}

func (l log) Error(msg string, args ...any) {
	l.ErrorContext(l.ctx, msg, args...)
}

func (l log) Errorf(format string, args ...any) {
	l.ErrorContext(l.ctx, fmt.Sprintf(format, args...))
}

func (l log) Trace(msg string, args ...any) {
	l.Log(l.ctx, LevelTrace, msg, args...)
}

func (l log) Tracef(format string, args ...any) {
	l.Log(l.ctx, LevelTrace, fmt.Sprintf(format, args...))
}

func (l log) Fatal(msg string, args ...any) {
	l.Log(l.ctx, LevelFatal, msg, args...)
}

func (l log) Fatalf(format string, args ...any) {
	l.Log(l.ctx, LevelFatal, fmt.Sprintf(format, args...))
}

func (l log) LogLevelf(ctx context.Context, level slog.Level, format string, args ...any) {
	l.Log(ctx, level, fmt.Sprintf(format, args...))
}

func (l log) With(args ...any) Logger {
	return log{
		l.Logger.With(args...),
		l.ctx,
	}
}

func (l log) LogLevel(ctx context.Context, level slog.Level, msg string, args ...any) {
	l.Log(ctx, level, msg, args...)
}

func FromContext(ctx context.Context) Logger {
	config.mu.RLock()
	defer config.mu.RUnlock()
	logger := &log{
		config.slogDefault,
		ctx,
	}
	return logger
}
