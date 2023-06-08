package mongodb

import (
	"context"
	"time"

	"github.com/lucasd-coder/user-manger-service/config"
	"github.com/lucasd-coder/user-manger-service/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var client *mongo.Client

func SetUpMongoDB(ctx context.Context, cfg *config.Config) {
	log := logger.FromContext(ctx)

	opts := options.Client().ApplyURI(cfg.MongoDB.URL)
	opts.Monitor = otelmongo.NewMonitor()
	opts.SetConnectTimeout(parseDuration(cfg.MongoDBConnTimeout))

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
	if err := client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func parseDuration(d string) time.Duration {
	const defaultDuration = 60 * time.Second

	pd, err := time.ParseDuration(d)
	if err != nil {
		pd = defaultDuration
	}
	return pd
}
