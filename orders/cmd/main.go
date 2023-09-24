package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"orders/internal/controller"
	"orders/pkg/config"
	"orders/pkg/handlers"
	"os"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	cfg, err := config.ReadConfig("config.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", cfg.Database.Host, cfg.Database.Port)))
	if err != nil {
		logger.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database(cfg.Database.Name)

	ctx := handlers.Context{
		Logger: logger,
		DB:     db,
	}

	router := chi.NewRouter()
	controller.RoutePaths(router, ctx)

	server := &http.Server{}
	server.Addr = fmt.Sprintf(":%d", cfg.Server.Port)
	server.Handler = router
	logger.Println(fmt.Sprintf("Server listen %d port", cfg.Server.Port))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
	logger.Println("Shutting down")
}
