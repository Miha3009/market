package controller

import (
	"encoding/json"
	"strconv"

	"github.com/miha3009/market/orders/internal/model"
	"github.com/miha3009/market/orders/internal/service"
	"github.com/miha3009/market/orders/pkg/handlers"

	"github.com/go-chi/chi/v5"
)

func HandleGetOrderById(ctx handlers.Context) (any, error) {
	id := chi.URLParam(ctx.R, "id")
	product, err := service.NewOrderService(ctx).SelectById(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func HandleGetOrders(ctx handlers.Context) (any, error) {
	userId, err := strconv.Atoi(ctx.R.URL.Query().Get("user"))
	if err != nil {
		return nil, err
	}

	return service.NewOrderService(ctx).SelectByUser(userId)
}

func HandleCreateOrder(ctx handlers.Context) (any, error) {
	var request model.Order
	decoder := json.NewDecoder(ctx.R.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, err
	}

	return service.NewOrderService(ctx).Create(request)
}

func HandleDeleteOrder(ctx handlers.Context) error {
	id := chi.URLParam(ctx.R, "id")
	return service.NewOrderService(ctx).Delete(id)
}

func RoutePaths(router chi.Router, ctx handlers.Context) {
	router.MethodFunc("GET", "/api/order/{id}", handlers.MakeJsonHandler(HandleGetOrderById, ctx))
	router.MethodFunc("GET", "/api/order", handlers.MakeJsonHandler(HandleGetOrders, ctx))
	router.MethodFunc("POST", "/api/order", handlers.MakeJsonHandler(HandleCreateOrder, ctx))
	router.MethodFunc("DELETE", "/api/order/{id}", handlers.MakeHandler(HandleDeleteOrder, ctx))
}
