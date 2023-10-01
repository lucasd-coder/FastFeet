package user

import (
	"time"

	"github.com/google/wire"
	"github.com/lucasd-coder/router-service/config"
	"github.com/lucasd-coder/router-service/internal/provider/publish"
	"github.com/lucasd-coder/router-service/internal/provider/validator"
	"github.com/lucasd-coder/router-service/internal/shared"
)

var InitializeService = wire.NewSet(
	wire.Bind(new(Service), new(*ServiceImpl)),
	NewService,
)

type ServiceImpl struct {
	validate     shared.Validator
	publish      shared.Publish
	cfg          *config.Config
	businessRepo shared.BusinessRepository
}

func NewService(
	val *validator.Validation,
	publish *publish.Published,
	cfg *config.Config,
	businessRepo shared.BusinessRepository,
) *ServiceImpl {
	return &ServiceImpl{
		validate:     val,
		publish:      publish,
		cfg:          cfg,
		businessRepo: businessRepo,
	}
}

func (s *ServiceImpl) getEventDate() string {
	return time.Now().Format(time.RFC3339)
}
