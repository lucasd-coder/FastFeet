package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lucasd-coder/fast-feet/order-data-service/config"
	model "github.com/lucasd-coder/fast-feet/order-data-service/internal/domain/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
	config     *config.Config
	connection *mongo.Client
}

func NewOrderRepository(cfg *config.Config, con *mongo.Client) *OrderRepository {
	return &OrderRepository{
		config:     cfg,
		connection: con,
	}
}

func (repo *OrderRepository) Save(ctx context.Context, order *model.Order) (*model.Order, error) {
	database := repo.connection.Database(repo.config.MongoDatabase)

	collection := repo.config.MongoCollections.Order.Collection

	result, err := database.Collection(collection).
		InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		ID: result.InsertedID.(primitive.ObjectID),
	}, nil
}

func (repo *OrderRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	database := repo.connection.Database(repo.config.MongoDatabase)

	collection := repo.config.MongoCollections.Order.Collection

	objectID, err := objectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objectID,
	}

	result := database.Collection(collection).FindOne(ctx, filter)

	return decode(result)
}

func (repo *OrderRepository) FindAll(ctx context.Context, pld *model.GetAllOrderRequest) ([]model.Order, error) {
	database := repo.connection.Database(repo.config.MongoDatabase)

	collection := repo.config.MongoCollections.Order.Collection

	filter, err := repo.extractFilterGetAllOrder(pld)
	if err != nil {
		return nil, fmt.Errorf("fail when extractFilter err: %w", err)
	}

	queryCtx, queryCancel := context.WithTimeout(ctx, repo.config.MongoCollections.Order.MaxTime)

	defer queryCancel()

	opt := options.Find()
	opt.SetSkip(pld.Offset)
	opt.SetLimit(pld.GetLimit())

	result, err := database.Collection(collection).Find(queryCtx, filter, opt)
	if err != nil {
		return nil, err
	}

	defer result.Close(ctx)

	var orders []model.Order
	order := new(model.Order)

	for result.Next(ctx) {
		if err := result.Decode(order); err != nil {
			return nil, fmt.Errorf("fail mongo cursor decode: %w", err)
		}
		orders = append(orders, *order)
	}

	return orders, nil
}

func decode(r *mongo.SingleResult) (*model.Order, error) {
	order := new(model.Order)
	if err := r.Decode(order); err != nil {
		return nil, fmt.Errorf("fail mongo decode: %w", err)
	}

	return order, nil
}

func objectIDFromHex(obj interface{}) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", obj))
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid object id: %w", err)
	}
	return objectID, nil
}

func (repo *OrderRepository) extractFilterGetAllOrder(pld *model.GetAllOrderRequest) (bson.M, error) {
	filter := bson.M{
		"deliverymanId": pld.DeliverymanID,
	}

	if pld.ID != "" {
		objectID, err := objectIDFromHex(pld.ID)
		if err != nil {
			return nil, err
		}

		filter["_id"] = objectID
	}

	fieldsTime := map[string]string{
		"startDate":  pld.StartDate,
		"endDate":    pld.EndDate,
		"createdAt":  pld.CreatedAt,
		"updatedAt":  pld.UpdatedAt,
		"canceledAt": pld.CanceledAt,
	}
	if err := repo.addTimeFilter(filter, fieldsTime); err != nil {
		return nil, err
	}

	fieldsString := map[string]string{
		"addresses.address":      pld.Address.Address,
		"addresses.postalCode":   pld.Address.PostalCode,
		"addresses.neighborhood": pld.Address.Neighborhood,
		"addresses.city":         pld.Address.City,
		"addresses.state":        pld.Address.State,
		"product.name":           pld.Product.Name,
	}
	repo.addStringFilter(filter, fieldsString)

	fieldsNumber := map[string]int64{
		"addresses.number": int64(pld.Address.Number),
	}

	repo.addNumberFilter(filter, fieldsNumber)

	return filter, nil
}

func (repo *OrderRepository) addStringFilter(filter primitive.M, fields map[string]string) {
	for fieldName, value := range fields {
		if value != "" {
			filter[fieldName] = value
		}
	}
}

func (repo *OrderRepository) addTimeFilter(filter primitive.M, fields map[string]string) error {
	for fieldName, value := range fields {
		if value != "" {
			parseTime, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return err
			}
			filter[fieldName] = bson.M{"$gte": primitive.NewDateTimeFromTime(parseTime)}
		}
	}
	return nil
}

func (repo *OrderRepository) addNumberFilter(filter primitive.M, fields map[string]int64) {
	for fieldName, value := range fields {
		if value != 0 {
			filter[fieldName] = value
		}
	}
}
