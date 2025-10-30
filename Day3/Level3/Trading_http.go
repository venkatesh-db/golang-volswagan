package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Order struct {
	ID     string  `json:"id"`
	Symbol string  `json:"symbol"`
	Side   string  `json:"side"`
	Price  float64 `json:"price"`
}

var orders = []Order{}

func main() {
	router := gin.Default()

	router.POST("/order", func(c *gin.Context) {
		var order Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		orders = append(orders, order)
		c.JSON(http.StatusOK, gin.H{"message": "Order placed", "data": order})
	})

	router.GET("/orders", func(c *gin.Context) {
		c.JSON(http.StatusOK, orders)
	})

	router.Run(":8081")
}

