package controller

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/miha3009/market/products/internal/model"
	"github.com/miha3009/market/products/internal/service"
	"github.com/miha3009/market/products/pkg/handlers"
)

func HandleGetProductById(ctx handlers.Context) (any, error) {
	id, err := strconv.Atoi(chi.URLParam(ctx.R, "id"))
	if err != nil {
		return nil, err
	}

	product, err := service.NewProductService(ctx).SelectById(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func HandleGetProducts(ctx handlers.Context) (any, error) {
	limit, err := strconv.Atoi(ctx.R.URL.Query().Get("limit"))
	if err != nil {
		ctx.Logger.Println(err)
		limit = math.MaxInt
	}

	offset, err := strconv.Atoi(ctx.R.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	return service.NewProductService(ctx).Select(offset, limit)
}

func HandleCreateProduct(ctx handlers.Context) (any, error) {
	var request model.Product
	decoder := json.NewDecoder(ctx.R.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return nil, err
	}

	return service.NewProductService(ctx).Create(request)
}

func HandleUpdateProduct(ctx handlers.Context) error {
	var request model.Product
	decoder := json.NewDecoder(ctx.R.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return err
	}

	return service.NewProductService(ctx).Update(request)
}

func HandleDeleteProduct(ctx handlers.Context) error {
	id, err := strconv.Atoi(chi.URLParam(ctx.R, "id"))
	if err != nil {
		return err
	}

	return service.NewProductService(ctx).Delete(id)
}

func RoutePaths(router chi.Router, ctx handlers.Context) {
	router.MethodFunc("GET", "/api/product/{id}", handlers.MakeJsonHandler(HandleGetProductById, ctx))
	router.MethodFunc("GET", "/api/product", handlers.MakeJsonHandler(HandleGetProducts, ctx))
	router.MethodFunc("POST", "/api/product", handlers.MakeJsonHandler(HandleCreateProduct, ctx))
	router.MethodFunc("PUT", "/api/product", handlers.MakeHandler(HandleUpdateProduct, ctx))
	router.MethodFunc("DELETE", "/api/product/{id}", handlers.MakeHandler(HandleDeleteProduct, ctx))
}
