package authservice

import (
	"context"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lucasd-coder/business-service/config"
	"golang.org/x/oauth2"
)

type options struct {
	transport        *http.Transport
	requestTimeout   time.Duration
	retryWaitTime    time.Duration
	retryMaxWaitTime time.Duration
	maxRetries       int
	url              string
	debug            bool
}

func NewClient(cfg *config.Config) *resty.Client {
	client := resty.New()

	opt := newOptions(cfg)

	client.EnableTrace().
		SetBaseURL(opt.url).
		SetRetryCount(cfg.MaxRetries).
		SetTransport(opt.transport).
		SetDebug(opt.debug).
		SetTimeout(opt.requestTimeout).
		SetRetryCount(opt.maxRetries).
		SetRetryMaxWaitTime(opt.retryMaxWaitTime).
		SetRetryWaitTime(opt.retryWaitTime)

	return client
}

func NewClientWithAuth(ctx context.Context, cfg *config.Config) (*resty.Client, error) {
	conf := &oauth2.Config{
		ClientID:     cfg.KeyCloakClientID,
		ClientSecret: cfg.KeyCloakClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL:  cfg.KeyCloakTokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
	}

	token, err := conf.PasswordCredentialsToken(ctx, cfg.KeyCloakUsername, cfg.KeyCloakPassword)
	if err != nil {
		return nil, err
	}

	clientConf := conf.Client(ctx, token)

	client := resty.NewWithClient(clientConf)

	opt := newOptions(cfg)

	client.EnableTrace().
		SetBaseURL(opt.url).
		SetDebug(opt.debug).
		SetTimeout(opt.requestTimeout)

	return client, err
}

func newOptions(cfg *config.Config) *options {
	transport := &http.Transport{
		MaxIdleConns:          cfg.AuthServiceMaxConn,
		IdleConnTimeout:       cfg.AuthServiceConnTimeout,
		MaxConnsPerHost:       cfg.AuthServiceMaxRoutes,
		ResponseHeaderTimeout: cfg.AuthServiceReadTimeout,
	}

	opt := &options{
		transport:        transport,
		requestTimeout:   cfg.AuthServiceRequestTimeout,
		retryWaitTime:    cfg.AuthServiceRetryWaitTime,
		retryMaxWaitTime: cfg.AuthServiceRetryMaxWaitTime,
		url:              cfg.AuthServiceURL,
		debug:            cfg.AuthServiceDebug,
		maxRetries:       cfg.AuthServiceMaxRetries,
	}

	return opt
}
