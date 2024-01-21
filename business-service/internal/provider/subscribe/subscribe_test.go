package subscribe_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lucasd-coder/fast-feet/business-service/internal/provider/subscribe"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/queueoptions"
	"github.com/lucasd-coder/fast-feet/pkg/monitor"
	testcontainers "github.com/lucasd-coder/fast-feet/pkg/testcontainers/rabbitmq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"gocloud.dev/pubsub"
)

type SubscribeSuite struct {
	suite.Suite
	ctx               context.Context
	rabbitmqContainer *rabbitmq.RabbitMQContainer
	opt               *queueoptions.Options
	queueName         string
	exchangeName      string
	mtr               monitor.Metrics
}

func (suite *SubscribeSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.queueName = "user-events-queue"
	suite.exchangeName = "user-events-exchange"
	var err error
	suite.rabbitmqContainer, err = testcontainers.RunContainer(suite.ctx, suite.queueName, suite.exchangeName)
	if err != nil {
		suite.T().Fatal(err)
	}
	host, err := suite.rabbitmqContainer.AmqpURL(suite.ctx)
	if err != nil {
		suite.T().Fatal(err)
	}
	// amqp://admin:admin123@localhost:5672/fastfeet
	os.Setenv("RABBIT_SERVER_URL", "amqp://guest:guest@"+removeUserAndPass(host)+"/fastfeet")

	slog.Info("get url ", slog.String("RABBIT_SERVER_URL", os.Getenv("RABBIT_SERVER_URL")))
}

func (suite *SubscribeSuite) TearDownSuite() {
	if err := suite.rabbitmqContainer.Terminate(suite.ctx); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *SubscribeSuite) SetupTest() {
	suite.opt = &queueoptions.Options{
		MaxConcurrentMessages:    1,
		MaxRetries:               2,
		WaitingTime:              time.Millisecond * 5,
		NumberOfMessageReceivers: 1,
		PollDelay:                time.Second * 2,
		MaxReceiveMessage:        time.Millisecond * 1,
		QueueURL:                 "rabbit://" + suite.queueName,
	}
	mtr, err := monitor.CreateMetrics(suite.queueName, prometheus.NewRegistry())
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.mtr = mtr
}

func (suite *SubscribeSuite) TestSubscribeConsumedMaxRetries() {
	count := 0
	var wg sync.WaitGroup

	wg.Add(2)

	suite.publish([]byte("erro"))

	handler := func(ctx context.Context, m []byte) error {
		count++
		defer wg.Done()
		if count == 2 {
			return nil
		}
		if string(m) == "erro" {
			return errors.New("error not found")
		}
		slog.Info("MSG: ", slog.Any("msg", m))
		return nil
	}
	sub := subscribe.New(handler, suite.opt, suite.mtr)

	go func() {
		sub.Start(suite.ctx)
	}()

	wg.Wait()

	suite.Equal(count, suite.opt.MaxRetries)
}
func TestSubscribeSuite(t *testing.T) {
	suite.Run(t, new(SubscribeSuite))
}

func removeUserAndPass(input string) string {
	index := strings.Index(input, "@")

	if index != -1 {
		return input[index+1:]
	}

	return input
}

func (suite *SubscribeSuite) publish(msg []byte) {
	client, err := pubsub.OpenTopic(suite.ctx, "rabbit://"+suite.exchangeName)
	if err != nil {
		suite.T().Fatal(err)
	}

	defer func() {
		if err := client.Shutdown(suite.ctx); err != nil {
			suite.T().Fatal(err)
		}
	}()

	m := pubsub.Message{
		Body: msg,
	}

	if err := client.Send(suite.ctx, &m); err != nil {
		suite.T().Fatal(err)
	}
}
