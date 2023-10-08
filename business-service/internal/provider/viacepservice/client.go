package viacepservice

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lucasd-coder/fast-feet/business-service/config"
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
		SetRetryCount(cfg.ViaCepMaxRetries).
		SetTransport(opt.transport).
		SetDebug(opt.debug).
		SetTimeout(opt.requestTimeout).
		SetRetryCount(opt.maxRetries).
		SetRetryMaxWaitTime(opt.retryMaxWaitTime).
		SetRetryWaitTime(opt.retryWaitTime)

	return client
}

func newOptions(cfg *config.Config) *options {
	transport := &http.Transport{
		MaxIdleConns:          cfg.ViaCepMaxConn,
		IdleConnTimeout:       cfg.ViaCepConnTimeout,
		MaxConnsPerHost:       cfg.ViaCepMaxRoutes,
		ResponseHeaderTimeout: cfg.ViaCepReadTimeout,
	}

	return &options{
		transport:        transport,
		requestTimeout:   cfg.ViaCepRequestTimeout,
		url:              cfg.ViaCepURL,
		debug:            cfg.ViaCepDebug,
		maxRetries:       cfg.ViaCepMaxRetries,
		retryWaitTime:    cfg.ViaCepRetryWaitTime,
		retryMaxWaitTime: cfg.ViaCepRetryMaxWaitTime,
	}
}
