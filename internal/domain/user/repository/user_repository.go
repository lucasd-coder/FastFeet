package repository

import (
	"context"
	"fmt"

	"github.com/lucasd-coder/user-manger-service/config"
	model "github.com/lucasd-coder/user-manger-service/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Config     *config.Config
	Connection *mongo.Client
}

func NewUserRepository(cfg *config.Config, con *mongo.Client) *UserRepository {
	return &UserRepository{
		Config:     cfg,
		Connection: con,
	}
}

func (repo *UserRepository) Save(ctx context.Context, user *model.User) (*model.User, error) {
	database := repo.Connection.Database(repo.Config.MongoDatabase)

	collection := repo.Config.MongoCollections.User.Collection

	result, err := database.Collection(collection).InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID: result.InsertedID.(primitive.ObjectID),
	}, nil
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	database := repo.Connection.Database(repo.Config.MongoDatabase)

	collection := repo.Config.MongoCollections.User.Collection

	filter := bson.M{
		"email": email,
	}

	result := database.Collection(collection).FindOne(ctx, filter)

	return decode(result)
}

func (repo *UserRepository) FindByUserID(ctx context.Context, userID string) (*model.User, error) {
	database := repo.Connection.Database(repo.Config.MongoDatabase)

	collection := repo.Config.MongoCollections.User.Collection

	filter := bson.M{
		"userId": userID,
	}

	result := database.Collection(collection).FindOne(ctx, filter)

	return decode(result)
}

func (repo *UserRepository) FindByCpf(ctx context.Context, cpf string) (*model.User, error) {
	database := repo.Connection.Database(repo.Config.MongoDatabase)

	collection := repo.Config.MongoCollections.User.Collection

	filter := bson.M{
		"cpf": cpf,
	}

	result := database.Collection(collection).FindOne(ctx, filter)

	return decode(result)
}

func decode(result *mongo.SingleResult) (*model.User, error) {
	user := new(model.User)
	if err := result.Decode(user); err != nil {
		return nil, fmt.Errorf("fail decode mongo result err: %w", err)
	}

	return user, nil
}
