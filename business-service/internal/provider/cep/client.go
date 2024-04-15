package cep

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

func NewClient(opt *options) *resty.Client {
	client := resty.New()

	client.EnableTrace().
		SetBaseURL(opt.url).
		SetTransport(otelhttp.NewTransport(opt.transport)).
		SetDebug(opt.debug).
		SetTimeout(opt.requestTimeout).
		SetRetryCount(opt.maxRetries).
		SetRetryMaxWaitTime(opt.retryMaxWaitTime).
		SetRetryWaitTime(opt.retryWaitTime)

	return client
}
