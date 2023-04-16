package logger

import (
	"context"
	"os"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/lucasd-coder/business-service/config"
	logrus "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Log struct {
	cfg *config.Config
}

func NewLog(cfg *config.Config) *Log {
	return &Log{cfg: cfg}
}

func FromContext(ctx context.Context) Logger {
	config := config.GetConfig()
	log := NewLog(config).GetGRPCLogger()

	logger := &logger{
		logger: log.WithContext(ctx),
	}

	return logger
}

func (l *Log) GetGRPCLogger() *logrus.Entry {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(l.cfg.ReportCaller)
	log.SetOutput(os.Stdout)

	logLevel, _ := logrus.ParseLevel(l.cfg.Level)
	logrus.SetLevel(logLevel)
	return log.WithFields(logrus.Fields{
		"logName":  l.cfg.App.Name,
		"logIndex": "message",
	})
}

func (l *Log) GetGRPCUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpc_logrus.UnaryServerInterceptor(l.GetGRPCLogger())
}

func (l *Log) GetGRPCStreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpc_logrus.StreamServerInterceptor(l.GetGRPCLogger())
}

type logger struct {
	logger *logrus.Entry
}

func (l *logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *logger) WithFields(keyValues map[string]interface{}) Logger {
	newEntry := l.logger.WithFields(convertToLogrusFields(keyValues))

	newLogger := &logger{
		logger: newEntry,
	}

	return newLogger
}

func (l *logger) WithField(key string, value interface{}) Logger {
	newEntry := l.logger.WithField(key, value)
	newLogger := &logger{
		logger: newEntry,
	}

	return newLogger
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, val := range fields {
		logrusFields[index] = val
	}
	return logrusFields
}
