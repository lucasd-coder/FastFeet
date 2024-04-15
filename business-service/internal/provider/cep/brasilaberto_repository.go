package cep

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"

	"github.com/go-resty/resty/v2"
	cacheProvider "github.com/lucasd-coder/fast-feet/business-service/internal/provider/cache"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/codec"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
)

var client *resty.Client

type BrasilAbertoRepository struct {
	cfg             *config.Config
	cacheRepository shared.CacheRepository[shared.AddressResponse]
}

type AddressResponse struct {
	Result struct {
		Street         string `json:"street"`
		Complement     string `json:"complement"`
		District       string `json:"district"`
		City           string `json:"city"`
		State          string `json:"state"`
		StateShortName string `json:"stateShortname"`
		ZipCode        string `json:"zipcode"`
	} `json:"result"`
}

func (a *AddressResponse) toAddress() *shared.AddressResponse {
	if a == nil {
		return nil
	}
	return &shared.AddressResponse{
		Address:      a.Result.Street,
		PostalCode:   a.Result.ZipCode,
		Neighborhood: a.Result.District,
		City:         a.Result.City,
		State:        a.Result.State,
	}
}

func NewBrasilAbertoRepository(cfg *config.Config,
	redisClient *redis.Client) *BrasilAbertoRepository {
	cacheRepository := cacheProvider.NewCacheRepository[shared.AddressResponse](redisClient)
	return &BrasilAbertoRepository{cfg, cacheRepository}
}

func (r *BrasilAbertoRepository) GetAddress(ctx context.Context, cep string) (*shared.AddressResponse, error) {
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

func (r *BrasilAbertoRepository) getAddress(ctx context.Context, cep string) (*shared.AddressResponse, error) {
	log := logger.FromContext(ctx)
	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))

	client := r.newClient()

	request := client.R().SetContext(ctx).SetLogger(log)

	response, err := request.
		SetPathParam("cep", cep).
		SetResult(&AddressResponse{}).
		SetError(&shared.HTTPError{}).
		Get("/v1/zipcode/{cep}")
	if err != nil {
		r.createSpanError(ctx, err, spanErrRequest)
		return nil, err
	}

	if response.IsError() {
		r.createSpanError(ctx, err, spanErrResponseStatus)
		return nil, fmt.Errorf(
			"err while execute request api-brasilaberto with statusCode: %s. Endpoint: /ws/{cep}/json, Method: GET", response.Status())
	}

	res, ok := response.Result().(*AddressResponse)
	if !ok {
		r.createSpanError(ctx, err, spanErrExtractResponse)
		return nil, fmt.Errorf("%w. Endpoint: /v1/zipcode/{cep}", shared.ErrExtractResponse)
	}

	log.Debugf("api-brasilaberto call successful. Endpoint: /v1/zipcode/{cep}, Method: GET, Response time: %s",
		response.ReceivedAt().String())

	return res.toAddress(), nil
}

func (r *BrasilAbertoRepository) getCachedToAddress(ctx context.Context, cep string) (*shared.AddressResponse, error) {
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

func (r *BrasilAbertoRepository) setCacheAndReturn(ctx context.Context, cep string) (*shared.AddressResponse, error) {
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

func (r *BrasilAbertoRepository) createSpanError(ctx context.Context, err error, msg string) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, msg)
	span.RecordError(err)
}

func (r *BrasilAbertoRepository) newOptions() *options {
	transport := &http.Transport{
		MaxIdleConns:          r.cfg.BrasilAbertoMaxConn,
		IdleConnTimeout:       r.cfg.BrasilAbertoConnTimeout,
		MaxConnsPerHost:       r.cfg.BrasilAbertoMaxRoutes,
		ResponseHeaderTimeout: r.cfg.BrasilAbertoReadTimeout,
	}

	return &options{
		transport:        transport,
		requestTimeout:   r.cfg.BrasilAbertoRequestTimeout,
		url:              r.cfg.BrasilAbertoURL,
		debug:            *r.cfg.BrasilAbertoDebug,
		maxRetries:       r.cfg.BrasilAbertoMaxRetries,
		retryWaitTime:    r.cfg.BrasilAbertoRetryWaitTime,
		retryMaxWaitTime: r.cfg.BrasilAbertoRetryMaxWaitTime,
	}
}

func (r *BrasilAbertoRepository) newClient() *resty.Client {
	if client == nil {
		client = NewClient(r.newOptions())
	}
	return client
}
