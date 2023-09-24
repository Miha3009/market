package repository

import (
	"context"
	"orders/internal/model"
	"orders/pkg/handlers"

	"go.mongodb.org/mongo-driver/bson"
)

func Select(ctx *handlers.Context, id int) (model.Order, error) {
	var result model.Order
	entry := ctx.DB.Collection("orders").FindOne(context.TODO(), bson.M{"id": id})
	err := entry.Decode(&result)
	return result, err
}

func InsertOne(ctx *handlers.Context, order model.Order) {
	/*order := model.Order{
		Id:   4,
		Date: time.Now(),
		Buyer: model.User{
			Id:   5,
			Name: "Someone",
		},
	}*/
	ctx.DB.Collection("orders").InsertOne(context.TODO(), order)
}
