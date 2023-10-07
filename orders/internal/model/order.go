package model

import (
	"time"

	pb "github.com/miha3009/market/protocol"
)

type Product struct {
	Id    int `json:"id"`
	Count int `json:"count"`
}

type Order struct {
	Id       string    `bson:"_id,omitempty" json:"id"`
	UserId   int       `bson:"userId" json:"userId"`
	Date     time.Time `bson:"date" json:"date"`
	Products []Product `bson:"products" json:"products"`
}

func ToReserveRequest(products []Product) []*pb.ReserveRequestProduct {
	res := make([]*pb.ReserveRequestProduct, len(products))
	for i := range products {
		res[i] = &pb.ReserveRequestProduct{Id: int32(products[i].Id), Count: int32(products[i].Count)}
	}
	return res
}
