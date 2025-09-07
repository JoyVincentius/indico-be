package handler

import (
	"fmt"
	"net/http"

	"indico-be/internal/models"
	"indico-be/internal/service"

	"github.com/gin-gonic/gin"
)

type orderRequest struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	BuyerID   string `json:"buyer_id" binding:"required"`
}

// RegisterOrderRoutes attaches order endpoints to the router.
func RegisterOrderRoutes(r *gin.Engine, svc *service.OrderService) {
	orders := r.Group("/orders")
	{
		orders.POST("", createOrder(svc))
		orders.GET("/:id", getOrder(svc))
	}
}

func createOrder(svc *service.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req orderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		order := &models.Order{
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			BuyerID:   req.BuyerID,
		}
		if err := svc.PlaceOrder(c.Request.Context(), order); err != nil {
			if err.Error() == "OUT_OF_STOCK" {
				c.JSON(http.StatusConflict, gin.H{"error": "OUT_OF_STOCK"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusCreated, order)
	}
}

func getOrder(svc *service.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		var id uint64
		if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		order, err := svc.GetOrder(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusOK, order)
	}
}
