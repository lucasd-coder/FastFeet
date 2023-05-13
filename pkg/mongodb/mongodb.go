package mongodb

import (
	"context"

	"github.com/lucasd-coder/order-data-service/config"
	"github.com/lucasd-coder/order-data-service/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func SetUpMongoDB(ctx context.Context, cfg *config.Config) {
	log := logger.FromContext(ctx)

	opts := options.Client().ApplyURI(cfg.MongoDB.URL)
	opts.SetConnectTimeout(cfg.MongoDBConnTimeout)

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
