package logger

type Fields map[string]interface{}

type Logger interface {
	Info(msg string, args ...any)

	Debug(msg string, args ...any)

	Error(msg string, args ...any)

	Warn(msg string, args ...any)

	Debugf(format string, args ...any)

	Infof(format string, args ...any)

	Warnf(format string, args ...any)

	Errorf(format string, args ...any)
}
