package http

import "github.com/gin-gonic/gin"

type Server struct {
	engine *gin.Engine
}

func NewServer() (*Server, error) {
	r := gin.Default()

	r.Use(corsMiddleware())
	n := &Server{
		engine: r,
	}

	n.initRouter()

	return n, nil
}

func (n *Server) Start(port string) error {
	return n.engine.Run(port)
}

func (n *Server) initRouter() {
	n.engine.GET("/ping")

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
