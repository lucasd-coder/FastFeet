package cep

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"

	"github.com/go-resty/resty/v2"
	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/codec"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ViaCepRepository struct {
	cfg             *config.Config
	cacheRepository shared.CacheRepository[shared.AddressResponse]
}

func (r *ViaCepRepository) GetAddress(ctx context.Context, cep string) (*shared.AddressResponse, error) {
	log := logger.FromContext(ctx)
	span := trace.SpanFromContext(ctx)
	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))

	address, err := r.getCachedToAddress(ctx, cep)
	if err != nil {
		span.AddEvent("setCacheAndReturn")
		span.SetAttributes(attribute.String("cep", cep))
		span.RecordError(err)
		log.Errorf("failed to retrieve cached address for cep with CEP: %s, err: %+v", cep, err)
		return r.setCacheAndReturn(ctx, cep)
	}
	span.AddEvent("getCachedToAddress")
	span.SetAttributes(attribute.String("cep", cep))
	log.Infof("successfully retrieved cached address for cep with CEP: %s", cep)

	return address, nil
}

func (r *ViaCepRepository) getAddress(ctx context.Context, cep string) (*shared.AddressResponse, error) {
	log := logger.FromContext(ctx)
	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))

	client := r.newClient()

	request := client.R().SetContext(ctx).SetLogger(log)

	response, err := request.
		SetPathParam("cep", cep).
		SetResult(&shared.AddressResponse{}).
		SetError(&shared.HTTPError{}).
		Get("/ws/{cep}/json/")
	if err != nil {
		r.createSpanError(ctx, err, spanErrRequest)
		return nil, err
	}

	if response.IsError() {
		r.createSpanError(ctx, err, spanErrResponseStatus)
		return nil, fmt.Errorf(
			"err while execute request api-viacep with statusCode: %s. Endpoint: /ws/{cep}/json, Method: GET", response.Status())
	}

	res, ok := response.Result().(*shared.AddressResponse)
	if !ok {
		r.createSpanError(ctx, err, spanErrExtractResponse)
		return nil, fmt.Errorf("%w. Endpoint: /ws/{cep}/json", shared.ErrExtractResponse)
	}

	log.Debugf("api-viacep call successful. Endpoint: /ws/{cep}/json, Method: GET, Response time: %s",
		response.ReceivedAt().String())

	return res, nil
}

func (r *ViaCepRepository) getCachedToAddress(ctx context.Context, cep string) (*shared.AddressResponse, error) {
	resultCache, err := r.cacheRepository.Get(ctx, cep)
	if err != nil {
		return nil, err
	}

	enc := codec.New[shared.AddressResponse]()

	var address *shared.AddressResponse

	if err := enc.Decode([]byte(resultCache), address); err != nil {
		return nil, err
	}

	return address, nil
}

func (r *ViaCepRepository) setCacheAndReturn(ctx context.Context, cep string) (*shared.AddressResponse, error) {
	log := logger.FromContext(ctx)

	address, err := r.getAddress(ctx, cep)
	if err != nil {
		return nil, err
	}

	if err := r.cacheRepository.Save(ctx, cep, *address, r.cfg.RedisTTL); err != nil {
		log.Errorf("fail with save cache repository with err: %+v", err)
		return address, nil
	}

	return address, nil
}

func (r *ViaCepRepository) createSpanError(ctx context.Context, err error, msg string) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, msg)
	span.RecordError(err)
}

func (r *ViaCepRepository) newOptions() *options {
	transport := &http.Transport{
		MaxIdleConns:          r.cfg.ViaCepMaxConn,
		IdleConnTimeout:       r.cfg.ViaCepConnTimeout,
		MaxConnsPerHost:       r.cfg.ViaCepMaxRoutes,
		ResponseHeaderTimeout: r.cfg.ViaCepReadTimeout,
	}

	return &options{
		transport:        transport,
		requestTimeout:   r.cfg.ViaCepRequestTimeout,
		url:              r.cfg.ViaCepURL,
		debug:            *r.cfg.ViaCepDebug,
		maxRetries:       r.cfg.ViaCepMaxRetries,
		retryWaitTime:    r.cfg.ViaCepRetryWaitTime,
		retryMaxWaitTime: r.cfg.ViaCepRetryMaxWaitTime,
	}
}

func (r *ViaCepRepository) newClient() *resty.Client {
	if client == nil {
		client = NewClient(r.newOptions())
	}
	return client
}
