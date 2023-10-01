package shared

import (
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/pkg/monitor"
	"github.com/lucasd-coder/order-data-service/config"
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
