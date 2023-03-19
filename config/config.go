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
		GrpcClient `env-required:"true" yaml:"grpc"`
		RabbitMQ   `env-required:"true" yaml:"rabbit-mq"`
	}

	GrpcClient struct {
		UserManagerService `env-required:"true" yaml:"user-manager-service"`
	}

	UserManagerService struct {
		URL      string `env-required:"true" yaml:"url"`
		MaxRetry uint   `env-required:"true" yaml:"max-retry"`
	}

	RabbitMQ struct {
		Queue `env-required:"true" yaml:"queue"`
	}

	Queue struct {
		QueueUserEvents `env-required:"true" yaml:"user-events"`
	}

	QueueUserEvents struct {
		URL                     string  `env-required:"true" yaml:"url"`
		MaxReceiveMessage       string  `yaml:"max-receive-message" env-default:"30s"`
		MaxRetries              int     `yaml:"max-retries" env-default:"5"`
		MaxConcurrentMessages   int     `yaml:"max-concurrent-messages" env-default:"10"`
		PollDelayInMilliseconds int     `yaml:"poll-delay-in-milliseconds" env-default:"100"`
		WaitingTime             float64 `yaml:"waiting-time" env-default:"2"`
	}
)

func ExportConfig(config *Config) {
	cfg = config
}

func GetConfig() *Config {
	return cfg
}
