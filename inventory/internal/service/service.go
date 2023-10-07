package service

import (
	"context"
	"database/sql"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/miha3009/market/inventory/internal/model"
	"github.com/miha3009/market/inventory/internal/repository"
	pb "github.com/miha3009/market/protocol"
)

type InventoryService struct {
	pb.UnimplementedInvetroryServer
	Logger *log.Logger
	DB     repository.InventoryRepository
}

func NewInventoryService(logger *log.Logger, db *sql.DB) *InventoryService {
	return &InventoryService{
		Logger: logger,
		DB:     repository.NewInventoryRepository(db),
	}
}

func (s *InventoryService) CheckAvaliable(c context.Context, r *pb.AvailabilityRequest) (*pb.AvailabilityResponse, error) {
	res, err := s.DB.Avaliable(r.Ids)
	return &pb.AvailabilityResponse{Available: res}, err
}

func (s *InventoryService) ToReserveRequest(r []*pb.ReserveRequestProduct) []model.ReserveRequest {
	res := make([]model.ReserveRequest, len(r))
	for i := range r {
		res[i].ProductId = int(r[i].GetId())
		res[i].Count = int(r[i].GetCount())
	}
	return res
}

func (s *InventoryService) Reserve(c context.Context, r *pb.ReserveRequest) (*pb.ReserveResponse, error) {
	succ, err := s.DB.Reserve(s.ToReserveRequest(r.Products))
	return &pb.ReserveResponse{Success: succ}, err
}

func (s *InventoryService) CancelReserve(c context.Context, r *pb.ReserveRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.DB.CancelReserve(s.ToReserveRequest(r.Products))
}
