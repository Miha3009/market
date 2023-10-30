package repository

import (
	"context"

	"github.com/miha3009/market/orders/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	SelectById(id string) (model.Order, error)
	SelectByUser(userId int) ([]model.Order, error)
	Create(value model.Order) (string, error)
	Delete(id string) error
}

type OrderRepositoryImpl struct {
	db *mongo.Database
}

func NewOrderRepository(db *mongo.Database) OrderRepository {
	return &OrderRepositoryImpl{
		db: db,
	}
}

func (r *OrderRepositoryImpl) SelectById(id string) (model.Order, error) {
	var result model.Order
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	entry := r.db.Collection("orders").FindOne(context.TODO(), bson.M{"_id": objId})
	err = entry.Decode(&result)
	return result, err
}

func (r *OrderRepositoryImpl) SelectByUser(userId int) ([]model.Order, error) {
	cur, err := r.db.Collection("orders").Find(context.TODO(), bson.M{"userId": userId})
	if err != nil {
		return nil, err
	}

	var result []model.Order
	if err = cur.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *OrderRepositoryImpl) Create(value model.Order) (string, error) {
	res, err := r.db.Collection("orders").InsertOne(context.TODO(), value)
	if err != nil {
		return "", err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", nil
}

func (r *OrderRepositoryImpl) Delete(id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.db.Collection("orders").DeleteOne(context.TODO(), bson.M{"_id": objId})
	return err
}
