package config

import "time"

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
		AesKey  string `env-required:"true" yaml:"aes-key" env:"AES_KEY"`
	}

	Log struct {
		Level        string `env-required:"true" yaml:"log-level"   env:"LOG_LEVEL"`
		ReportCaller bool   `yaml:"report-caller" default:"false"`
	}

	GRPC struct {
		Port string `env-required:"true" yaml:"port" env:"GRPC_PORT"`
	}

	Integration struct {
		GrpcClient `env-required:"true" yaml:"grpc"`
		HTTPClint  `env-required:"true" yaml:"http"`
		RabbitMQ   `env-required:"true" yaml:"rabbit-mq"`
		Redis      `env-required:"true" yaml:"redis"`
	}

	GrpcClient struct {
		UserManagerService `env-required:"true" yaml:"user-manager-service"`
		OrderDataService   `env-required:"true" yaml:"order-data-service"`
	}

	UserManagerService struct {
		URL      string `env-required:"true" yaml:"url"`
		MaxRetry uint   `env-required:"true" yaml:"max-retry"`
	}

	OrderDataService struct {
		URL      string `env-required:"true" yaml:"url"`
		MaxRetry uint   `env-required:"true" yaml:"max-retry"`
	}

	RabbitMQ struct {
		Queue `env-required:"true" yaml:"queue"`
	}

	Queue struct {
		QueueUserEvents  `env-required:"true" yaml:"user-events"`
		QueueOrderEvents `env-required:"true" yaml:"order-events"`
	}

	HTTPClint struct {
		AccessAuthService `env-required:"true" yaml:"access-auth-service"`
		KeyCloak          `env-required:"true" yaml:"keycloak"`
		ViaCep            `env-required:"true" yaml:"viacep"`
	}

	QueueUserEvents struct {
		URL                      string        `env-required:"true" yaml:"url"`
		MaxReceiveMessage        time.Duration `yaml:"max-receive-message" env-default:"30s"`
		MaxRetries               int           `yaml:"max-retries" env-default:"5"`
		MaxConcurrentMessages    int           `yaml:"max-concurrent-messages" env-default:"10"`
		NumberOfMessageReceivers int           `yaml:"number-of-message-receivers" env-default:"2"`
		PollDelay                time.Duration `yaml:"poll-delay-in-milliseconds" env-default:"100ms"`
		WaitingTime              time.Duration `yaml:"waiting-time" env-default:"2s"`
	}

	QueueOrderEvents struct {
		URL                      string        `env-required:"true" yaml:"url"`
		MaxReceiveMessage        time.Duration `yaml:"max-receive-message" env-default:"30s"`
		MaxRetries               int           `yaml:"max-retries" env-default:"5"`
		MaxConcurrentMessages    int           `yaml:"max-concurrent-messages" env-default:"10"`
		NumberOfMessageReceivers int           `yaml:"number-of-message-receivers" env-default:"2"`
		PollDelay                time.Duration `yaml:"poll-delay-in-milliseconds" env-default:"100ms"`
		WaitingTime              time.Duration `yaml:"waiting-time" env-default:"2s"`
	}

	AccessAuthService struct {
		AuthServiceURL              string        `env-required:"true" yaml:"url"`
		AuthServiceMaxConn          int           `env-required:"true" yaml:"max-conn"`
		AuthServiceMaxRoutes        int           `env-required:"true" yaml:"max-routes"`
		AuthServiceReadTimeout      time.Duration `yaml:"read-timeout" env-default:"60s"`
		AuthServiceConnTimeout      time.Duration `yaml:"conn-timeout" env-default:"60s"`
		AuthServiceDebug            bool          `yaml:"debug" env-default:"true"`
		AuthServiceRequestTimeout   time.Duration `env-required:"true" yaml:"request-timeout"`
		AuthServiceMaxRetries       int           `env-required:"true" yaml:"max-retry"`
		AuthServiceRetryWaitTime    time.Duration `env-required:"true" yaml:"retry-wait-time"`
		AuthServiceRetryMaxWaitTime time.Duration `env-required:"true" yaml:"retry-max-wait-time"`
	}

	KeyCloak struct {
		KeyCloakTokenURL       string        `env-required:"true" yaml:"token-url"`
		KeyCloakUsername       string        `env-required:"true" yaml:"username"`
		KeyCloakPassword       string        `env-required:"true" yaml:"password"`
		KeyCloakRequestTimeout time.Duration `env-required:"true" yaml:"request-timeout"`
		KeyCloakClientID       string        `env-required:"true" yaml:"client-id"`
		KeyCloakClientSecret   string        `env-required:"true" yaml:"client-secret"`
	}

	ViaCep struct {
		ViaCepURL              string        `env-required:"true" yaml:"url"`
		ViaCepMaxConn          int           `env-required:"true" yaml:"max-conn"`
		ViaCepMaxRoutes        int           `env-required:"true" yaml:"max-routes"`
		ViaCepReadTimeout      time.Duration `yaml:"read-timeout" env-default:"60s"`
		ViaCepConnTimeout      time.Duration `yaml:"conn-timeout" env-default:"60s"`
		ViaCepDebug            bool          `yaml:"debug" env-default:"true"`
		ViaCepRequestTimeout   time.Duration `env-required:"true" yaml:"request-timeout"`
		ViaCepMaxRetries       int           `env-required:"true" yaml:"max-retry"`
		ViaCepRetryWaitTime    time.Duration `env-required:"true" yaml:"retry-wait-time"`
		ViaCepRetryMaxWaitTime time.Duration `env-required:"true" yaml:"retry-max-wait-time"`
	}

	Redis struct {
		RedisURL      string        `env-required:"true" yaml:"url" env:"REDIS_URL"`
		RedisDB       int           `env-required:"true" yaml:"db" env:"REDIS_DB"`
		RedisPassword string        `env-required:"true" yaml:"password" env:"REDIS_HOST_PASSWORD"`
		RedisTTL      time.Duration `env-required:"true" yaml:"ttl"`
	}
)

func ExportConfig(config *Config) {
	cfg = config
}

func GetConfig() *Config {
	return cfg
}
