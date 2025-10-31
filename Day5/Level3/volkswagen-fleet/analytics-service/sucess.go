
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Telemetry struct {
	VehicleID   string    `json:"vehicle_id"`
	Speed       float64   `json:"speed"`
	Temperature float64   `json:"temperature"`
	FuelLevel   float64   `json:"fuel_level"`
	Timestamp   time.Time `json:"timestamp"`
}

var telemetryData []Telemetry

// Compute basic scoring
func computeDriverScore(t Telemetry) int {
	score := 100
	if t.Speed > 120 {
		score -= 20
	}
	if t.FuelLevel < 15 {
		score -= 10
	}
	if t.Temperature > 90 {
		score -= 15
	}
	return score
}

func main() {
	r := gin.Default()

	r.GET("/analytics/:vehicleID", func(c *gin.Context) {
		vid := c.Param("vehicleID")
		var latest Telemetry
		found := false
		for i := len(telemetryData) - 1; i >= 0; i-- {
			if telemetryData[i].VehicleID == vid {
				latest = telemetryData[i]
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle data not found"})
			return
		}
		score := computeDriverScore(latest)
		c.JSON(http.StatusOK, gin.H{
			"vehicle_id": latest.VehicleID,
			"score":      score,
			"status":     "computed",
		})
	})

	log.Println("Fleet Analytics Service running on :8081")
	r.Run(":8081")
}
