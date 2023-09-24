package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	Logger *log.Logger
	DB     *sql.DB
}

type ContextHandlerFunc func(Context) error

type ContextHandlerFuncWithResponse func(Context) (any, error)

func MakeHandler(handler ContextHandlerFunc, ctx Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, r *http.Request) {
		err := handler(Context{
			W:      responseWriter,
			R:      r,
			Logger: ctx.Logger,
			DB:     ctx.DB,
		})

		if err != nil {
			ctx.Logger.Println(err)
			responseWriter.WriteHeader(http.StatusNotFound)
		}
	}
}

func MakeJsonHandler(handler ContextHandlerFuncWithResponse, ctx Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, r *http.Request) {
		obj, err := handler(Context{
			W:      responseWriter,
			R:      r,
			Logger: ctx.Logger,
			DB:     ctx.DB,
		})

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
