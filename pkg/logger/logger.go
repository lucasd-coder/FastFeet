package logger

import (
	"context"
	"os"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/lucasd-coder/business-service/config"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Log struct {
	cfg *config.Config
}

func NewLog(cfg *config.Config) *Log {
	return &Log{cfg: cfg}
}

func FromContext(ctx context.Context) *logger.Entry {
	config := config.GetConfig()
	log := NewLog(config).GetGRPCLogger()
	return log.WithContext(ctx)
}

func (l *Log) GetGRPCLogger() *logger.Entry {
	log := logger.New()
	log.SetFormatter(&logger.JSONFormatter{})
	log.SetOutput(os.Stdout)

	logLevel, _ := logger.ParseLevel(l.cfg.Level)
	logger.SetLevel(logLevel)
	return log.WithFields(logger.Fields{
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
