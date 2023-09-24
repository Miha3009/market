package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"products/internal/controller"
	"products/pkg/config"
	"products/pkg/handlers"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/miha3009/market/protocol"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	cfg, err := config.ReadConfig("config.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	client, err := sql.Open("postgres", connString)
	if err != nil {
		logger.Fatal(err)
	}
	defer client.Close()

	inventoryConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Inventory.Host, cfg.Inventory.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}
	defer inventoryConn.Close()
	pb.

	ctx := handlers.Context{
		Logger: logger,
		DB:     client,
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
