package order_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	noProviderVal "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"github.com/lucasd-coder/fast-feet/router-service/config"
	"github.com/lucasd-coder/fast-feet/router-service/internal/domain/order"
	"github.com/lucasd-coder/fast-feet/router-service/internal/mocks"
	"github.com/lucasd-coder/fast-feet/router-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/router-service/internal/shared"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateOrderSuite struct {
	suite.Suite
	cfg          config.Config
	id           string
	objectId     primitive.ObjectID
	svc          order.ServiceImpl
	businessRepo *mocks.BusinessRepository_internal_shared
	publishRepo  *mocks.Publish_internal_shared
	ctx          context.Context
	valErrs      noProviderVal.ValidationErrors
}

func (suite *CreateOrderSuite) SetupSuite() {
	baseDir, err := os.Getwd()
	if err != nil {
		suite.T().Errorf("os.Getwd() error = %v", err)
		return
	}

	staticDir := filepath.Join(baseDir, "..", "..", "..", "/config/config-test.yml")

	slog.Info("config lod", "dir", staticDir)
	err = cleanenv.ReadConfig(staticDir, &suite.cfg)
	if err != nil {
		suite.T().Errorf("cleanenv.ReadConfig() error = %v", err)
		return
	}
	config.ExportConfig(&suite.cfg)
	optlogger := shared.NewOptLogger(&suite.cfg)
	logger := logger.NewLogger(optlogger)
	logDefault := logger.GetLog()
	slog.SetDefault(logDefault)
	suite.ctx = context.Background()
}

func (suite *CreateOrderSuite) SetupTest() {
	val := validator.NewValidation()

	suite.businessRepo = new(mocks.BusinessRepository_internal_shared)
	suite.publishRepo = new(mocks.Publish_internal_shared)
	suite.id = uuid.NewString()

	suite.objectId = primitive.NewObjectID()

	suite.svc = *order.NewService(val, suite.publishRepo, &suite.cfg, suite.businessRepo)
}

func (suite *CreateOrderSuite) TearDownTest() {
	suite.businessRepo.AssertExpectations(suite.T())
	suite.publishRepo.AssertExpectations(suite.T())
	suite.businessRepo = nil
	suite.publishRepo = nil
}

func (suite *CreateOrderSuite) TestSave() {
	tests := []struct {
		name  string
		args  *order.Order
		mocks func()
		check func(name string, err error)
	}{
		{
			name: "test validation field userID",
			args: &order.Order{
				DeliverymanID: suite.id,
				Product: order.Product{
					Name: "mesa",
				},
				Address: order.Address{
					PostalCode: "123456",
					Number:     20,
				},
			},
			check: func(name string, err error) {
				if err != nil {
					valErrs := suite.valErrs
					suite.True(errors.As(err, &valErrs))
				}
			},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if tt.mocks != nil {
				tt.mocks()
			}
			err := suite.svc.Save(suite.ctx, tt.args)
			if tt.check != nil {
				tt.check(tt.name, err)
			}
		})
	}
}

func TestCreateOrderSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderSuite))
}
