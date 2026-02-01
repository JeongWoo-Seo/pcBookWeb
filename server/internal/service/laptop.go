package service

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

type LaptopService struct {
	rdb *redis.Client
}

func NewLaptopService(rdb *redis.Client) *LaptopService {
	return &LaptopService{rdb: rdb}
}

func (s *LaptopService) GetActiveLaptopList(
	ctx context.Context,
) ([]string, error) {

	keys, err := s.rdb.Keys(ctx, "laptop:alive:*").Result()
	if err != nil {
		return nil, err
	}

	laptops := make([]string, 0, len(keys))
	for _, key := range keys {
		id := strings.TrimPrefix(key, "laptop:alive:")
		laptops = append(laptops, id)
	}

	return laptops, nil
}
