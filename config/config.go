package config

import "time"

var cfg *Config

type (
	Config struct {
		App    `yaml:"app"`
		Server `yaml:"server"`
		Log    `yaml:"logger"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	Server struct {
		Port         string        `env-required:"true" yaml:"port" env:"SERVER_PORT"`
		ReadTimeout  time.Duration `yaml:"readTimeout" default:"10s"`
		WriteTimeout time.Duration `yaml:"writeTimeout" default:"10s"`
	}
)

func ExportConfig(config *Config) {
	cfg = config
}

func GetConfig() *Config {
	return cfg
}
