package order

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/wire"
	"github.com/lucasd-coder/business-service/internal/shared"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
)

var InitializeService = wire.NewSet(
	wire.Bind(new(Service), new(*ServiceImpl)),
	NewService,
)

type ServiceImpl struct {
	validate         shared.Validator
	orderRepository  Repository
	authRepository   shared.AuthRepository
	viaCepRepository ViaCepRepository
}

func NewService(
	val shared.Validator,
	orderRepo Repository,
	authRepo shared.AuthRepository,
	viaCepRepo ViaCepRepository,
) *ServiceImpl {
	return &ServiceImpl{
		validate:         val,
		orderRepository:  orderRepo,
		authRepository:   authRepo,
		viaCepRepository: viaCepRepo,
	}
}

func (s *ServiceImpl) hasActiveUser(ctx context.Context, id string) error {
	log := logger.FromContext(ctx)

	log.Infof("get started to check is active user with id: %s", id)

	isActiveUser, err := s.authRepository.IsActiveUser(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrUserNotFound) {
			return shared.NotFoundError(shared.ErrUserNotFound)
		}
		return err
	}

	if !isActiveUser.Active {
		log.Errorf("deliveryman not active with id: %s", id)
		return fmt.Errorf("%w: deliveryman not active with id: %s", shared.ErrUserUnauthorized, id)
	}
	return nil
}

func (s *ServiceImpl) hasPermissionIsAdmin(ctx context.Context, id string) (bool, error) {
	log := logger.FromContext(ctx)

	log.Infof("get started roles with id: %s", id)

	roles, err := s.authRepository.FindRolesByID(ctx, id)
	if err != nil {
		log.Errorf("error when check permission with id: %s, err: %v", id, err)
		return false, err
	}

	for _, role := range roles.Roles {
		if strings.EqualFold(shared.ADMIN, role) {
			return true, nil
		}
	}

	return false, nil
}
