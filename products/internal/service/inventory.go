package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	pb "github.com/miha3009/market/protocol"
)

const cacheTTL = 30 * time.Second

type InventoryService interface {
	Check(id int) bool
	CheckRange(ids []int) []bool
}

type InventoryServiceImpl struct {
	client pb.InvetroryClient
	cache  *redis.Client
	logger *log.Logger
}

func NewInvetroryService(client pb.InvetroryClient, cache *redis.Client, logger *log.Logger) InventoryService {
	return &InventoryServiceImpl{
		client: client,
		cache:  cache,
		logger: logger,
	}
}

func (s *InventoryServiceImpl) Check(id int) bool {
	return s.CheckRange([]int{id})[0]
}

func (s *InventoryServiceImpl) CheckRange(ids []int) []bool {
	result := make([]bool, len(ids))
	notFoundIds := make([]int32, 0)
	notFoundIdsIdx := make([]int, 0)
	for i := range ids {
		value, err := s.cache.Get(context.TODO(), strconv.Itoa(ids[i])).Result()
		if err == redis.Nil {
			notFoundIds = append(notFoundIds, int32(ids[i]))
			notFoundIdsIdx = append(notFoundIdsIdx, i)
			continue
		} else if err != nil {
			s.logger.Println(err)
			return result
		}
		avaliable, err := strconv.ParseBool(value)
		if err != nil {
			s.logger.Println(err)
			return result
		}
		result[i] = avaliable
	}

	if len(notFoundIds) > 0 {
		resp, err := s.client.CheckAvaliable(context.TODO(), &pb.AvailabilityRequest{Ids: notFoundIds})
		if err != nil {
			s.logger.Println(err)
			return result
		}
		for i := range notFoundIds {
			result[notFoundIdsIdx[i]] = resp.Available[i]
			s.cache.Set(context.TODO(), strconv.Itoa(int(notFoundIds[i])), strconv.FormatBool(resp.Available[i]), cacheTTL)
		}
	}

	return result
}
