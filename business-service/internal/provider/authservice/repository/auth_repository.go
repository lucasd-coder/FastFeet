package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/lucasd-coder/business-service/config"
	"github.com/lucasd-coder/business-service/internal/provider/authservice"
	"github.com/lucasd-coder/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	spanErrMarchal         = "Error json.Marshal"
	spanErrRequest         = "Request Error"
	spanErrResponseStatus  = "Response Status Error"
	spanErrExtractUserID   = "Error ExtractUserID"
	spanErrNewClient       = "Error NewClientWithAuth"
	spanErrExtractResponse = "Error Extract Response"
)

type AuthRepository struct {
	cfg *config.Config
}

func NewAuthRepository(cfg *config.Config) *AuthRepository {
	return &AuthRepository{cfg}
}

func (r *AuthRepository) Register(ctx context.Context, pld *shared.Register) (*shared.RegisterUserResponse, error) {
	log := logger.FromContext(ctx)
	span := trace.SpanFromContext(ctx)
	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))

	client := authservice.NewClient(r.cfg)

	request := client.R().SetContext(ctx)

	body, err := json.Marshal(pld)
	if err != nil {
		errMsg := fmt.Errorf("err while marshalling payload register: %w", err)
		r.createSpanError(ctx, err, spanErrMarchal)
		return nil, errMsg
	}

	response, err := request.SetBody(body).
		SetHeader("Content-Type", "application/json").
		SetError(&shared.HTTPError{}).
		SetResult(&shared.RegisterUserResponse{}).
		Post("/api/register")

	if err != nil {
		r.createSpanError(ctx, err, spanErrRequest)
		return nil, err
	}

	if response.IsError() {
		errMsg := fmt.Errorf(
			"err while execute request auth-service with statusCode: %s. Endpoint: /api/register, Method: POST", response.Status())
		r.createSpanError(ctx, err, spanErrResponseStatus)
		return nil, errMsg
	}

	res, err := r.extractUserID(response)
	if err != nil {
		r.createSpanError(ctx, err, spanErrExtractUserID)
		return nil, err
	}

	msg := fmt.Sprintf("auth-service call successful. Endpoint: /api/register, Method: POST, Response time: %s",
		response.ReceivedAt().String())
	span.SetStatus(codes.Ok, msg)
	log.Debug(msg)

	return res, err
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*shared.GetUserResponse, error) {
	log := logger.FromContext(ctx)
	span := trace.SpanFromContext(ctx)
	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))

	client, err := authservice.NewClientWithAuth(ctx, r.cfg)
	if err != nil {
		r.createSpanError(ctx, err, spanErrNewClient)
		return nil, err
	}

	request := client.R().SetContext(ctx)

	response, err := request.
		SetPathParam("email", email).
		SetResult(&shared.GetUserResponse{}).
		SetError(&shared.HTTPError{}).
		Get("/api/users/{email}")
	if err != nil {
		r.createSpanError(ctx, err, spanErrRequest)
		return nil, err
	}

	if response.IsError() {
		err, ok := response.Error().(*shared.HTTPError)
		r.createSpanError(ctx, err, spanErrResponseStatus)
		if ok {
			if strings.EqualFold(err.Message, shared.ErrUserNotFound.Error()) {
				return nil, shared.ErrUserNotFound
			}

			return nil, err
		}
		return nil, fmt.Errorf(
			"err while execute request auth-service with statusCode: %s. Endpoint: /api/users/{email}, Method: GET", response.Status())
	}

	res, ok := response.Result().(*shared.GetUserResponse)

	if !ok {
		errMsg := fmt.Errorf("%w. Endpoint: /api/users", shared.ErrExtractResponse)
		r.createSpanError(ctx, err, spanErrExtractResponse)
		return nil, errMsg
	}

	msg := fmt.Sprintf("auth-service call successful. Endpoint: /api/users, Method: GET, Response time: %s",
		response.ReceivedAt().String())

	log.Debug(msg)
	span.SetStatus(codes.Ok, msg)

	return res, nil
}

func (r *AuthRepository) FindRolesByID(ctx context.Context, id string) (*shared.GetRolesResponse, error) {
	log := logger.FromContext(ctx)
	span := trace.SpanFromContext(ctx)
	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))

	client, err := authservice.NewClientWithAuth(ctx, r.cfg)
	if err != nil {
		r.createSpanError(ctx, err, spanErrNewClient)
		return nil, err
	}

	request := client.R().SetContext(ctx)

	response, err := request.
		SetPathParam("id", id).
		SetResult(&shared.GetRolesResponse{}).
		SetError(&shared.HTTPError{}).
		Get("/api/users/roles/{id}")
	if err != nil {
		r.createSpanError(ctx, err, spanErrRequest)
		return nil, err
	}

	if response.IsError() {
		errMsg := fmt.Errorf(
			"err while execute request auth-service with statusCode: %s. Endpoint: /api/users/roles/{id}, Method: GET", response.Status())
		r.createSpanError(ctx, err, spanErrResponseStatus)
		return nil, errMsg
	}

	res, ok := response.Result().(*shared.GetRolesResponse)

	if !ok {
		errMsg := fmt.Errorf("%w. Endpoint: /api/users/roles/{id}", shared.ErrExtractResponse)
		r.createSpanError(ctx, err, spanErrExtractResponse)
		return nil, errMsg
	}

	msg := fmt.Sprintf("auth-service call successful. Endpoint: /api/users/roles/{id}, Method: GET, Response time: %s",
		response.ReceivedAt().String())

	log.Debug(msg)
	span.SetStatus(codes.Ok, msg)

	return res, nil
}

func (r *AuthRepository) IsActiveUser(ctx context.Context, id string) (*shared.IsActiveUser, error) {
	log := logger.FromContext(ctx)
	span := trace.SpanFromContext(ctx)
	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))

	client, err := authservice.NewClientWithAuth(ctx, r.cfg)
	if err != nil {
		r.createSpanError(ctx, err, spanErrNewClient)
		return nil, err
	}

	request := client.R().SetContext(ctx)

	response, err := request.
		SetPathParam("id", id).
		SetResult(&shared.IsActiveUser{}).
		SetError(&shared.HTTPError{}).
		Get("/api/users/is-active/{id}")
	if err != nil {
		r.createSpanError(ctx, err, spanErrRequest)
		return nil, err
	}

	if response.IsError() {
		if response.StatusCode() == http.StatusNotFound {
			span.RecordError(shared.ErrUserNotFound)
			return nil, shared.ErrUserNotFound
		}
		errMsg := fmt.Errorf(
			"err while execute request auth-service with statusCode: %s. Endpoint: /api/users/roles/{id}, Method: GET", response.Status())
		r.createSpanError(ctx, err, spanErrResponseStatus)
		return nil, errMsg
	}

	res, ok := response.Result().(*shared.IsActiveUser)

	if !ok {
		errMsg := fmt.Errorf("%w. Endpoint: /api/users/roles/{id}", shared.ErrExtractResponse)
		r.createSpanError(ctx, err, spanErrExtractResponse)
		return nil, errMsg
	}

	msg := fmt.Sprintf("auth-service call successful. Endpoint: /api/users/is-active/{id}, Method: GET, Response time: %s",
		response.ReceivedAt().String())

	log.Debugf(msg)
	span.SetStatus(codes.Ok, msg)

	return res, nil
}

func (r *AuthRepository) extractUserID(response *resty.Response) (*shared.RegisterUserResponse, error) {
	if response == nil {
		return nil, nil
	}

	location := response.Header().Get("Location")

	u, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("err whiling url parse extract location extractUserID: %w", err)
	}

	path := strings.Split(u.Path, "/")
	userID := path[len(path)-1]

	return &shared.RegisterUserResponse{
		ID: userID,
	}, nil
}

func (r *AuthRepository) createSpanError(ctx context.Context, err error, msg string) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, msg)
	span.RecordError(err)
}
