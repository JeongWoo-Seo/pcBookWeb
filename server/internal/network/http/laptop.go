package http

import (
	"net/http"

	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/service"
	"github.com/gin-gonic/gin"
)

type LaptopRouter struct {
	router        *HttpNetwork
	laptopService *service.LaptopService
}

func NewLaptopRouter(httpNetwork *HttpNetwork, laptopService *service.LaptopService) *LaptopRouter {
	r := &LaptopRouter{
		router:        httpNetwork,
		laptopService: laptopService,
	}

	api := httpNetwork.engine.Group("/laptop")
	{
		api.GET("/list", r.List)
	}

	return r
}

func (h *LaptopRouter) List(c *gin.Context) {
	laptops, err := h.laptopService.GetActiveLaptopList(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": laptops,
	})
}
