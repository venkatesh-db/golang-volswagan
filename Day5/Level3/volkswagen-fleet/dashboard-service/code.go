
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var vehicleAnalytics = []map[string]interface{}{
	{"vehicle_id": "VW001", "score": 85, "status": "Good"},
	{"vehicle_id": "VW002", "score": 70, "status": "Attention"},
}

func main() {
	r := gin.Default()

	r.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, vehicleAnalytics)
	})

	log.Println("Dashboard Service running on :8083")
	r.Run(":8083")
}
