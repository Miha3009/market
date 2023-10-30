package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"github.com/miha3009/market/inventory/internal/service"
	"github.com/miha3009/market/inventory/pkg/config"
	pb "github.com/miha3009/market/protocol"
	"google.golang.org/grpc"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	cfg, err := config.ReadConfig(os.Getenv("CONFIG_PATH"))
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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		logger.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterInvetroryServer(s, service.NewInventoryService(logger, client))
	logger.Println(fmt.Sprintf("Server listen %d port", cfg.Server.Port))
	if err := s.Serve(lis); err != nil {
		logger.Fatal(err)
	}

	logger.Println("Shutting down")
}
