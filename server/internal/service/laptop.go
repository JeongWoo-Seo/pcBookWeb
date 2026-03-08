package service

import (
	"context"
	"strconv"
	"time"

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

	expireAfter := int64(10) //10초 내 활성화 된 laptop
	now := time.Now().Unix()

	laptops, err := s.rdb.ZRangeByScore(
		ctx,
		"laptop:alive",
		&redis.ZRangeBy{
			Min: strconv.FormatInt(now-expireAfter, 10),
			Max: "+inf",
		},
	).Result()

	if err != nil {
		return nil, err
	}

	return laptops, nil
}
