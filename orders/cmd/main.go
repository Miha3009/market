package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/miha3009/market/orders/internal/controller"
	"github.com/miha3009/market/orders/pkg/config"
	"github.com/miha3009/market/orders/pkg/handlers"
	pb "github.com/miha3009/market/protocol"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

	inventoryConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Inventory.Host, cfg.Inventory.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}
	defer inventoryConn.Close()

	w := &kafka.Writer{
		Addr:  kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic: cfg.Kafka.Topic,
	}

	ctx := handlers.Context{
		Logger:    logger,
		DB:        db,
		Inventory: pb.NewInvetroryClient(inventoryConn),
		Kafka:     w,
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
