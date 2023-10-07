package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/miha3009/market/products/internal/controller"
	"github.com/miha3009/market/products/pkg/config"
	"github.com/miha3009/market/products/pkg/handlers"
	pb "github.com/miha3009/market/protocol"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	cfg, err := config.ReadConfig("config.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Name)
	postgresClient, err := sql.Open("postgres", connString)
	if err != nil {
		logger.Fatal(err)
	}
	defer postgresClient.Close()

	inventoryConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Inventory.Host, cfg.Inventory.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}
	defer inventoryConn.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	ctx := handlers.Context{
		Logger:    logger,
		DB:        postgresClient,
		Cache:     redisClient,
		Inventory: pb.NewInvetroryClient(inventoryConn),
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
