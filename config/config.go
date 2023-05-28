package config

import "time"

var cfg *Config

type (
	Config struct {
		App         `yaml:"app"`
		Server      `yaml:"server"`
		Log         `yaml:"logger"`
		Integration `yaml:"integration"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		AesKey  string `env-required:"true" yaml:"aes-key" env:"AES_KEY"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	Integration struct {
		RabbitMQ   `env-required:"true" yaml:"rabbit-mq"`
		GrpcClient `env-required:"true" yaml:"grpc"`
	}

	RabbitMQ struct {
		Topic `env-required:"true" yaml:"topic"`
	}

	Topic struct {
		TopicUserEvents  `env-required:"true" yaml:"user-events"`
		TopicOrderEvents `env-required:"true" yaml:"order-events"`
	}

	Server struct {
		Port         string        `env-required:"true" yaml:"port" env:"SERVER_PORT"`
		ReadTimeout  time.Duration `yaml:"readTimeout" default:"10s"`
		WriteTimeout time.Duration `yaml:"writeTimeout" default:"10s"`
	}

	TopicUserEvents struct {
		URL         string        `env-required:"true" yaml:"url"`
		MaxRetries  int           `yaml:"max-retries" env-default:"3"`
		WaitingTime time.Duration `yaml:"waiting-time" env-default:"2s"`
	}

	TopicOrderEvents struct {
		URL         string        `env-required:"true" yaml:"url"`
		MaxRetries  int           `yaml:"max-retries" env-default:"3"`
		WaitingTime time.Duration `yaml:"waiting-time" env-default:"2s"`
	}

	GrpcClient struct {
		BusinessService `env-required:"true" yaml:"business-service"`
	}

	BusinessService struct {
		URL      string `env-required:"true" yaml:"url"`
		MaxRetry uint   `env-required:"true" yaml:"max-retry"`
	}
)

func ExportConfig(config *Config) {
	cfg = config
}

func GetConfig() *Config {
	return cfg
}
