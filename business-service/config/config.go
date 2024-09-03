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
		AesKey  string `env-required:"true" yaml:"aes-key" env:"AES_KEY"`
	}

	Log struct {
		Level        string `env-required:"true" yaml:"log-level"   env:"LOG_LEVEL"`
		ReportCaller *bool  `yaml:"report-caller" default:"false"`
	}

	GRPC struct {
		Port string `env-required:"true" yaml:"port" env:"GRPC_PORT"`
	}

	HTTP struct {
		Port    string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		Timeout time.Duration `env-required:"true" yaml:"timeout"`
	}

	Integration struct {
		GrpcClient    `env-required:"true" yaml:"grpc"`
		HTTPClint     `env-required:"true" yaml:"http"`
		RabbitMQ      `env-required:"true" yaml:"rabbit-mq"`
		Redis         `env-required:"true" yaml:"redis"`
		OpenTelemetry `env-required:"true" yaml:"otlp"`
	}

	GrpcClient struct {
		UserManagerService `env-required:"true" yaml:"user-manager-service"`
		OrderDataService   `env-required:"true" yaml:"order-data-service"`
		AuthService        `env-required:"true" yaml:"auth-service"`
	}

	UserManagerService struct {
		URL                                string        `env-required:"true" yaml:"url" env:"USER_MANAGER_URL"`
		MaxRetry                           uint          `env-required:"true" yaml:"max-retry"`
		UserManagerServiceRetryWaitTime    time.Duration `env-required:"true" yaml:"retry-wait-time"`
		UserManagerServiceRetryMaxWaitTime time.Duration `env-required:"true" yaml:"retry-max-wait-time"`
	}

	OrderDataService struct {
		URL                              string        `env-required:"true" yaml:"url" env:"ORDER_DATA_URL"`
		MaxRetry                         uint          `env-required:"true" yaml:"max-retry"`
		OrderDataServiceRetryWaitTime    time.Duration `env-required:"true" yaml:"retry-wait-time"`
		OrderDataServiceRetryMaxWaitTime time.Duration `env-required:"true" yaml:"retry-max-wait-time"`
	}

	AuthService struct {
		URL                         string        `env-required:"true" yaml:"url" env:"AUTH_URL"`
		MaxRetry                    uint          `env-required:"true" yaml:"max-retry"`
		AuthServiceRetryWaitTime    time.Duration `env-required:"true" yaml:"retry-wait-time"`
		AuthServiceRetryMaxWaitTime time.Duration `env-required:"true" yaml:"retry-max-wait-time"`
	}

	RabbitMQ struct {
		Queue           `env-required:"true" yaml:"queue"`
		RabbitServerURL string `env-required:"true" env:"RABBIT_SERVER_URL"`
	}

	Queue struct {
		QueueUserEvents  `env-required:"true" yaml:"user-events"`
		QueueOrderEvents `env-required:"true" yaml:"order-events"`
	}

	HTTPClint struct {
		ViaCep       `env-required:"true" yaml:"viacep"`
		BrasilAberto `env-required:"true" yaml:"brasilaberto"`
	}

	QueueUserEvents struct {
		QueueURL                 string        `env-required:"true" yaml:"url"`
		MaxReceiveMessage        time.Duration `yaml:"max-receive-message" env-default:"60s"`
		MaxRetries               int           `yaml:"max-retries" env-default:"5"`
		MaxConcurrentMessages    int           `yaml:"max-concurrent-messages" env-default:"10"`
		NumberOfMessageReceivers int           `yaml:"number-of-message-receivers" env-default:"2"`
		PollDelay                time.Duration `yaml:"poll-delay-in-milliseconds" env-default:"100ms"`
		WaitingTime              time.Duration `yaml:"waiting-time" env-default:"2s"`
	}

	QueueOrderEvents struct {
		QueueURL                 string        `env-required:"true" yaml:"url"`
		MaxReceiveMessage        time.Duration `yaml:"max-receive-message" env-default:"60s"`
		MaxRetries               int           `yaml:"max-retries" env-default:"5"`
		MaxConcurrentMessages    int           `yaml:"max-concurrent-messages" env-default:"10"`
		NumberOfMessageReceivers int           `yaml:"number-of-message-receivers" env-default:"2"`
		PollDelay                time.Duration `yaml:"poll-delay-in-milliseconds" env-default:"100ms"`
		WaitingTime              time.Duration `yaml:"waiting-time" env-default:"2s"`
	}

	ViaCep struct {
		ViaCepURL              string        `env-required:"true" yaml:"url" env:"VIA_CEP_URL"`
		ViaCepMaxConn          int           `env-required:"true" yaml:"max-conn"`
		ViaCepMaxRoutes        int           `env-required:"true" yaml:"max-routes"`
		ViaCepReadTimeout      time.Duration `yaml:"read-timeout" env-default:"60s"`
		ViaCepConnTimeout      time.Duration `yaml:"conn-timeout" env-default:"60s"`
		ViaCepDebug            *bool         `yaml:"debug" env-default:"true"`
		ViaCepRequestTimeout   time.Duration `env-required:"true" yaml:"request-timeout"`
		ViaCepMaxRetries       int           `env-required:"true" yaml:"max-retry"`
		ViaCepRetryWaitTime    time.Duration `env-required:"true" yaml:"retry-wait-time"`
		ViaCepRetryMaxWaitTime time.Duration `env-required:"true" yaml:"retry-max-wait-time"`
		ViaCepEnabled          *bool         `yaml:"enabled" env-required:"true"`
	}

	BrasilAberto struct {
		BrasilAbertoURL              string        `env-required:"true" yaml:"url" env:"BRASIL_ABERTO_URL"`
		BrasilAbertoMaxConn          int           `env-required:"true" yaml:"max-conn"`
		BrasilAbertoMaxRoutes        int           `env-required:"true" yaml:"max-routes"`
		BrasilAbertoReadTimeout      time.Duration `yaml:"read-timeout" env-default:"60s"`
		BrasilAbertoConnTimeout      time.Duration `yaml:"conn-timeout" env-default:"60s"`
		BrasilAbertoDebug            *bool         `yaml:"debug" env-default:"true"`
		BrasilAbertoRequestTimeout   time.Duration `env-required:"true" yaml:"request-timeout"`
		BrasilAbertoMaxRetries       int           `env-required:"true" yaml:"max-retry"`
		BrasilAbertoRetryWaitTime    time.Duration `env-required:"true" yaml:"retry-wait-time"`
		BrasilAbertoRetryMaxWaitTime time.Duration `env-required:"true" yaml:"retry-max-wait-time"`
		BrasilAbertoEnabled          *bool         `env-required:"true" yaml:"enabled"`
	}
	Redis struct {
		RedisURL      string        `env-required:"true" yaml:"url" env:"REDIS_URL"`
		RedisDB       int           `env-required:"true" yaml:"db" env:"REDIS_DB"`
		RedisPassword string        `env-required:"true" yaml:"password" env:"REDIS_HOST_PASSWORD"`
		RedisTTL      time.Duration `env-required:"true" yaml:"ttl"`
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
