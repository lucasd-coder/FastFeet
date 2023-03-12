package config

var cfg *Config

type (
	Config struct {
		App         `yaml:"app"`
		GRPC        `yaml:"grpc"`
		Log         `yaml:"logger"`
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

	Integration struct {
		UserManagerService UserManagerService `env-required:"true" yaml:"user-manager-service"`
	}

	UserManagerService struct {
		URL      string `env-required:"true" yaml:"url"`
		MaxRetry int    `env-required:"true" yaml:"max-retry"`
	}
)

func ExportConfig(config *Config) {
	cfg = config
}

func GetConfig() *Config {
	return cfg
}
