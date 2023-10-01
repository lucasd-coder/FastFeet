package mongodb

import (
	"context"
	"time"

	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var client *mongo.Client

type Option struct {
	ConnTimeout time.Duration
	URL         string
}

func SetUpMongoDB(ctx context.Context, opt *Option) {
	log := logger.FromContext(ctx)

	opts := options.Client().ApplyURI(opt.URL)
	opts.Monitor = otelmongo.NewMonitor()
	opts.SetConnectTimeout(opt.ConnTimeout)

	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error MongoDB connection: %+v", err.Error())
		return
	}

	log.Info("MongoDB Connected")

	client = mongoClient
}

func GetClientMongoDB() *mongo.Client {
	return client
}

func CloseConnMongoDB(ctx context.Context) error {
	return client.Disconnect(ctx)
}
