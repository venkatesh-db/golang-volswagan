
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Telemetry represents a single vehicle data point
type Telemetry struct {
	VehicleID   string  `json:"vehicle_id"`
	Speed       float64 `json:"speed"`
	Temperature float64 `json:"temperature"`
	FuelLevel   float64 `json:"fuel_level"`
	Timestamp   time.Time `json:"timestamp"`
}

var telemetryData []Telemetry

func main() {
	r := gin.Default()

	// Receive telemetry data
	r.POST("/telemetry", func(c *gin.Context) {
		var t Telemetry
		if err := c.BindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		t.Timestamp = time.Now()
		telemetryData = append(telemetryData, t)
		c.JSON(http.StatusOK, gin.H{"status": "received"})
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "vehicle-telemetry-service up"})
	})

	log.Println("Vehicle Telemetry Service running on :8080")
	r.Run(":8080")
}
