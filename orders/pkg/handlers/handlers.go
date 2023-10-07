package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	pb "github.com/miha3009/market/protocol"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Context struct {
	W         http.ResponseWriter
	R         *http.Request
	Logger    *log.Logger
	DB        *mongo.Database
	Inventory pb.InvetroryClient
	Kafka     *kafka.Writer
}

type ContextHandlerFunc func(Context) error

type ContextHandlerFuncWithResponse func(Context) (any, error)

func NewContext(responseWriter http.ResponseWriter, r *http.Request, ctx Context) Context {
	return Context{
		W:         responseWriter,
		R:         r,
		Logger:    ctx.Logger,
		DB:        ctx.DB,
		Inventory: ctx.Inventory,
		Kafka:     ctx.Kafka,
	}
}

func MakeHandler(handler ContextHandlerFunc, ctx Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, r *http.Request) {
		err := handler(NewContext(responseWriter, r, ctx))

		if err != nil {
			ctx.Logger.Println(err)
			responseWriter.WriteHeader(http.StatusNotFound)
		}
	}
}

func MakeJsonHandler(handler ContextHandlerFuncWithResponse, ctx Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, r *http.Request) {
		obj, err := handler(NewContext(responseWriter, r, ctx))

		if err != nil || obj == nil {
			ctx.Logger.Println(err)
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}

		data, err := json.Marshal(obj)
		if err != nil {
			ctx.Logger.Println(err)
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Write(data)
	}
}
