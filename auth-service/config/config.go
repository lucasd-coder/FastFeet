package config

import "time"

var cfg *Config

type (
	Config struct {
		App         `yaml:"app"`
		GRPC        `yaml:"grpc"`
		HTTP        `yaml:"http"`
		Log         `yaml:"logger"`
		Integration `yaml:"integration"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	Log struct {
		Level        string `env-required:"true" yaml:"log-level"   env:"LOG_LEVEL"`
		ReportCaller bool   `yaml:"report-caller" default:"false"`
	}

	GRPC struct {
		Port string `env-required:"true" yaml:"port" env:"GRPC_PORT"`
	}

	HTTP struct {
		Port    string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		Timeout time.Duration `env-required:"true" yaml:"timeout"`
	}

	Integration struct {
		HTTPClint     `env-required:"true" yaml:"http"`
		OpenTelemetry `env-required:"true" yaml:"otlp"`
	}

	HTTPClint struct {
		KeyCloak `env-required:"true" yaml:"keycloak"`
	}

	KeyCloak struct {
		KeyCloakBaseURL        string        `env-required:"true" yaml:"base-url" env:"KEYCLOAK_URL"`
		KeyCloakUsername       string        `env-required:"true" yaml:"username" env:"KEYCLOAK_USERNAME"`
		KeyCloakPassword       string        `env-required:"true" yaml:"password" env:"KEYCLOAK_PASSWORD"`
		KeyCloakRealm          string        `env-required:"true" yaml:"realm" env:"KEYCLOAK_REALM"`
		KeyCloakRequestTimeout time.Duration `env-required:"true" yaml:"request-timeout"`
		KeyCloakDebug          bool          `yaml:"debug" env-default:"true"`
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
