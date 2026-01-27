package http

import (
	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/network/ws"
	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/service"
	"github.com/gin-gonic/gin"
)

type HttpNetwork struct {
	engine  *gin.Engine
	service *service.Service
}

func NewHttpNetwork(service *service.Service) (*HttpNetwork, error) {
	httpNetwork := &HttpNetwork{
		engine:  gin.Default(),
		service: service,
	}

	httpNetwork.engine.Use(corsMiddleware())

	hub := ws.NewHub()
	go hub.Run()
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
