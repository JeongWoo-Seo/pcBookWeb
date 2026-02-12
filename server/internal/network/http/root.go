package http

import (
	"context"

	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/network/ws"
	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type HttpNetwork struct {
	engine  *gin.Engine
	service *service.Service
	rdb     *redis.Client
}

func NewHttpNetwork(service *service.Service, rdb *redis.Client) (*HttpNetwork, error) {
	httpNetwork := &HttpNetwork{
		engine:  gin.Default(),
		service: service,
		rdb:     rdb,
	}

	httpNetwork.engine.Use(corsMiddleware())

	hub := ws.NewHub(rdb)
	go hub.Run(context.Background())
	httpNetwork.engine.GET("/ws", ws.HandleWebSocket(hub))

	NewLaptopRouter(httpNetwork, service.LaptopService)

	return httpNetwork, nil
}

func (n *HttpNetwork) Start(port string) error {
	return n.engine.Run(port)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
