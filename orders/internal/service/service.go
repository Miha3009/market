package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/miha3009/market/orders/internal/model"
	"github.com/miha3009/market/orders/internal/repository"
	"github.com/miha3009/market/orders/pkg/handlers"
	pb "github.com/miha3009/market/protocol"
	"github.com/segmentio/kafka-go"
)

type OrderService interface {
	SelectById(id string) (model.Order, error)
	SelectByUser(userId int) ([]model.Order, error)
	Create(value model.Order) (model.CreateResponse, error)
	Delete(id string) error
}

type OrderServiceImpl struct {
	repo     repository.OrderRepository
	inv      pb.InvetroryClient
	messageQ *kafka.Writer
}

func NewOrderService(ctx handlers.Context) OrderService {
	return &OrderServiceImpl{
		repo:     repository.NewOrderRepository(ctx.DB),
		inv:      ctx.Inventory,
		messageQ: ctx.Kafka,
	}
}

func (s *OrderServiceImpl) SelectById(id string) (model.Order, error) {
	return s.repo.SelectById(id)
}

func (s *OrderServiceImpl) SelectByUser(userId int) ([]model.Order, error) {
	return s.repo.SelectByUser(userId)
}

func (s *OrderServiceImpl) Create(value model.Order) (model.CreateResponse, error) {
	resp, err := s.inv.Reserve(context.TODO(), &pb.ReserveRequest{Products: model.ToReserveRequest(value.Products)})
	if err != nil {
		return model.CreateResponse{}, err
	}
	if resp.Success == false {
		return model.CreateResponse{Succ: false}, nil
	}

	value.Date = time.Now()
	id, err := s.repo.Create(value)
	s.messageQ.WriteMessages(context.TODO(),
		kafka.Message{
			Key:   []byte(fmt.Sprintf("user%d@mail.ru", value.UserId)),
			Value: []byte(s.FormatMessage(value)),
		},
	)
	return model.CreateResponse{Id: id, Succ: true}, err
}

func (s *OrderServiceImpl) Delete(id string) error {
	order, err := s.SelectById(id)
	if err != nil {
		return err
	}

	_, err = s.inv.CancelReserve(context.TODO(), &pb.ReserveRequest{Products: model.ToReserveRequest(order.Products)})
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *OrderServiceImpl) FormatMessage(order model.Order) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Hello, User%d!\r\n", order.UserId))
	builder.WriteString("New order ")
	builder.WriteString(order.Date.Format(time.RFC1123))
	builder.WriteString(":\r\n")
	for i := range order.Products {
		builder.WriteString(fmt.Sprintf("Product%d, Count: %d\r\n", order.Products[i].Id, order.Products[i].Count))
	}
	return builder.String()
}
