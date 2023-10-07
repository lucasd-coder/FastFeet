package repository

import (
	"context"
	"fmt"
	"net/http/httptrace"

	"github.com/lucasd-coder/business-service/config"
	cacheProvider "github.com/lucasd-coder/business-service/internal/provider/cache"
	"github.com/lucasd-coder/business-service/internal/provider/viacepservice"
	"github.com/lucasd-coder/business-service/internal/shared"
	"github.com/lucasd-coder/business-service/internal/shared/codec"
	"github.com/lucasd-coder/business-service/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	spanErrMarchal         = "Error json.Marshal"
	spanErrRequest         = "Request Error"
	spanErrResponseStatus  = "Response Status Error"
	spanErrExtractResponse = "Error Extract Response"
)

type ViaCepRepository struct {
	cfg             *config.Config
	cacheRepository shared.CacheRepository[shared.ViaCepAddressResponse]
}

func NewViaCepRepository(cfg *config.Config,
	client *redis.Client) *ViaCepRepository {
	cacheRepository := cacheProvider.NewCacheRepository[shared.ViaCepAddressResponse](client)
	return &ViaCepRepository{cfg, cacheRepository}
}

func (r *ViaCepRepository) GetAddress(ctx context.Context, cep string) (*shared.ViaCepAddressResponse, error) {
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

func (r *ViaCepRepository) getAddress(ctx context.Context, cep string) (*shared.ViaCepAddressResponse, error) {
	log := logger.FromContext(ctx)
	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))

	client := viacepservice.NewClient(r.cfg)

	request := client.R().SetContext(ctx)

	response, err := request.
		SetPathParam("cep", cep).
		SetResult(&shared.ViaCepAddressResponse{}).
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

	res, ok := response.Result().(*shared.ViaCepAddressResponse)
	if !ok {
		r.createSpanError(ctx, err, spanErrExtractResponse)
		return nil, fmt.Errorf("%w. Endpoint: /ws/{cep}/json", shared.ErrExtractResponse)
	}

	log.Debugf("api-viacep call successful. Endpoint: /api/users, Method: GET, Response time: %s",
		response.ReceivedAt().String())

	return res, nil
}

func (r *ViaCepRepository) getCachedToAddress(ctx context.Context, cep string) (*shared.ViaCepAddressResponse, error) {
	resultCache, err := r.cacheRepository.Get(ctx, cep)
	if err != nil {
		return nil, err
	}

	enc := codec.New[shared.ViaCepAddressResponse]()

	var address *shared.ViaCepAddressResponse

	if err := enc.Decode([]byte(resultCache), address); err != nil {
		return nil, err
	}

	return address, nil
}

func (r *ViaCepRepository) setCacheAndReturn(ctx context.Context, cep string) (*shared.ViaCepAddressResponse, error) {
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
