package service

import "github.com/redis/go-redis/v9"

type Service struct {
	LaptopService *LaptopService
}

func NewService(rdb *redis.Client) *Service {
	return &Service{
		NewLaptopService(rdb),
	}
}
