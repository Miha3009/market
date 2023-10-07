package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	pb "github.com/miha3009/market/protocol"
)

type Context struct {
	W         http.ResponseWriter
	R         *http.Request
	Logger    *log.Logger
	DB        *sql.DB
	Cache     *redis.Client
	Inventory pb.InvetroryClient
}

type ContextHandlerFunc func(Context) error

type ContextHandlerFuncWithResponse func(Context) (any, error)

func NewContext(responseWriter http.ResponseWriter, r *http.Request, ctx Context) Context {
	return Context{
		W:         responseWriter,
		R:         r,
		Logger:    ctx.Logger,
		DB:        ctx.DB,
		Cache:     ctx.Cache,
		Inventory: ctx.Inventory,
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
