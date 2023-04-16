package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/lucasd-coder/business-service/config"
	model "github.com/lucasd-coder/business-service/internal/domain/user"
	"github.com/lucasd-coder/business-service/internal/provider/authservice"
	"github.com/lucasd-coder/business-service/internal/shared"
	"github.com/lucasd-coder/business-service/pkg/logger"
)

type AuthRepository struct {
	cfg *config.Config
}

func NewAuthRepository(cfg *config.Config) *AuthRepository {
	return &AuthRepository{cfg}
}

func (r *AuthRepository) Register(ctx context.Context, pld *model.Register) (*model.RegisterUserResponse, error) {
	log := logger.FromContext(ctx)

	client := authservice.NewClient(r.cfg)

	request := client.R()

	body, err := json.Marshal(pld)
	if err != nil {
		return nil, fmt.Errorf("err while marshalling payload register: %w", err)
	}

	response, err := request.SetBody(body).
		SetHeader("Content-Type", "application/json").
		SetError(&shared.HTTPError{}).
		Post("/api/register")
	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, fmt.Errorf(
			"err while execute request auth-service with statusCode: %s. Endpoint: /api/register, Method: POST", response.Status())
	}

	res, err := r.extractUserID(response)
	if err != nil {
		return nil, err
	}

	log.Debugf("auth-service call successful. Endpoint: /api/register, Method: POST, Response time: %s",
		response.ReceivedAt().String())

	return res, err
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*model.GetUserResponse, error) {
	log := logger.FromContext(ctx)

	client, err := authservice.NewClientWithAuth(ctx, r.cfg)
	if err != nil {
		return nil, err
	}

	request := client.R()

	response, err := request.
		SetPathParam("email", email).
		SetResult(&model.GetUserResponse{}).
		SetError(&shared.HTTPError{}).
		Get("/api/users/{email}")
	if err != nil {
		return nil, err
	}

	if response.IsError() {
		err, ok := response.Error().(*shared.HTTPError)

		if ok {
			if strings.EqualFold(err.Message, shared.ErrUserNotFound.Error()) {
				return nil, shared.ErrUserNotFound
			}

			return nil, err
		}

		return nil, fmt.Errorf(
			"err while execute request auth-service with statusCode: %s. Endpoint: /api/users, Method: GET", response.Status())
	}

	res, ok := response.Result().(*model.GetUserResponse)

	if !ok {
		return nil, fmt.Errorf("%w. Endpoint: /api/users", shared.ErrExtractResponse)
	}

	log.Debugf("auth-service call successful. Endpoint: /api/register, Method: POST, Response time: %s",
		response.ReceivedAt().String())

	return res, nil
}

func (r *AuthRepository) extractUserID(response *resty.Response) (*model.RegisterUserResponse, error) {
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

	return &model.RegisterUserResponse{
		ID: userID,
	}, nil
}
