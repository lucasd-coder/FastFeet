package config

import "time"

var cfg *Config

type (
	Config struct {
		App         `yaml:"app"`
		GRPC        `yaml:"grpc"`
		HTTP        `yaml:"http"`
		Log         `yaml:"logger"`
		MongoDB     `yaml:"mongodb"`
		Integration `yaml:"integration"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	GRPC struct {
		Port string `env-required:"true" yaml:"port" env:"GRPC_PORT"`
	}

	HTTP struct {
		Port    string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		Timeout time.Duration `env-required:"true" yaml:"timeout"`
	}

	MongoDB struct {
		URL                string           `env-required:"true" yaml:"url" env:"MONGODB_URL"`
		MongoDBConnTimeout time.Duration    `yaml:"connTimeout" default:"10s"`
		MongoDatabase      string           `env-required:"true" yaml:"database"`
		MongoCollections   MongoCollections `env-required:"true" yaml:"collections"`
	}
	MongoCollections struct {
		User `env-required:"true" yaml:"user"`
	}

	User struct {
		Collection string `env-required:"true" yaml:"collection"`
	}
	Integration struct {
		OpenTelemetry `env-required:"true" yaml:"otlp"`
	}

	OpenTelemetry struct {
		URL      string        `env-required:"true" yaml:"url" env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
		Protocol string        `env-required:"true" yaml:"protocol" env:"OTEL_EXPORTER_OTLP_PROTOCOL"`
		Timeout  time.Duration `env-required:"true" yaml:"timeout" env:"OTEL_EXPORTER_OTLP_TIMEOUT"`
	}
)

func ExportConfig(config *Config) {
	cfg = config
}

func GetConfig() *Config {
	return cfg
}
