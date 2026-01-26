package http

import (
	"net/http"

	"github.com/JeongWoo-Seo/pcBookWeb/server/internal/service"
	"github.com/gin-gonic/gin"
)

type LaptopRouter struct {
	laptopService *service.LaptopService
}

func NewLaptopRouter(laptopService *service.LaptopService) *LaptopRouter {
	return &LaptopRouter{
		laptopService: laptopService,
	}
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
