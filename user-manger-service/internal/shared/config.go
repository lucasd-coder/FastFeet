package shared

import (
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/pkg/mongodb"
	"github.com/lucasd-coder/fast-feet/pkg/monitor"
	"github.com/lucasd-coder/fast-feet/user-manger-service/config"
)

func NewOptLogger(cfg *config.Config) logger.Option {
	return logger.Option{
		AppName: cfg.Name,
		Level:   cfg.Level,
	}
}

func NewOptOtel(cfg *config.Config) monitor.Option {
	return monitor.Option{
		ServiceName: cfg.Name,
		Protocol:    cfg.OpenTelemetry.Protocol,
		URL:         cfg.OpenTelemetry.URL,
	}
}

func NewOptMongoDB(cfg *config.Config) mongodb.Option {
	return mongodb.Option{
		ConnTimeout: cfg.MongoDBConnTimeout,
		URL:         cfg.MongoDB.URL,
	}
}
