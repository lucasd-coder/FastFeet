package shared

import (
	"time"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/pkg/monitor"
	"github.com/lucasd-coder/router-service/config"
)

type Message struct {
	Body     []byte
	Metadata map[string]string
}

type CreateEvent struct {
	Message string `json:"message,omitempty"`
}

type Options struct {
	TopicURL    string
	MaxRetries  int
	WaitingTime time.Duration
}

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
