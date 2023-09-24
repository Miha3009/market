package controller

import (
	"encoding/json"
	"net/http"
	"orders/internal/repository"
	"orders/pkg/handlers"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Response struct {
	X string `json:"x"`
}

func HandleGetOrderById(ctx handlers.Context) {
	id, err := strconv.Atoi(chi.URLParam(ctx.R, "id"))
	if err != nil {
		ctx.Logger.Println(err)
		ctx.W.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := repository.Select(&ctx, id)
	if err != nil {
		ctx.Logger.Println(err)
		ctx.W.WriteHeader(http.StatusNotFound)
		return
	}

	data, _ := json.Marshal(result)
	ctx.W.Write(data)
	ctx.W.Header().Set("Content-Type", "application/json")
}

func RoutePaths(router chi.Router, ctx handlers.Context) {
	router.MethodFunc("GET", "/api/order/{id}", handlers.MakeHandler(HandleGetOrderById, ctx))
}
