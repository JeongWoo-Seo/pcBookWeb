package service

import "github.com/gin-gonic/gin"

type LaptopService struct {
	//rdb *redis.Client
}

func NewLaptopService() *LaptopService {
	//return &LaptopService{rdb: rdb}
	return &LaptopService{}
}

// func (s *LaptopService) GetActiveLaptopList(ctx *gin.Context) ([]string, error) {
// 	// Redis SET: active:laptops
// 	laptops, err := s.rdb.SMembers(ctx, "active:laptops").Result()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return laptops, nil
// }

func (s *LaptopService) GetActiveLaptopList(ctx *gin.Context) ([]string, error) {
	laptops := []string{
		"test1",
		"test2",
	}

	return laptops, nil
}
