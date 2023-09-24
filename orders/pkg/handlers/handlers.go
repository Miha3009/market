package handlers

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	Logger *log.Logger
	DB     *mongo.Database
}

type ContextHandlerFunc func(Context)

func MakeHandler(handler ContextHandlerFunc, ctx Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, r *http.Request) {
		handler(Context{
			W:      responseWriter,
			R:      r,
			Logger: ctx.Logger,
			DB:     ctx.DB,
		})
	}
}
