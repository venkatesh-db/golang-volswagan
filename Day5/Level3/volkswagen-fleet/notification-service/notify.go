package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type Alert struct {
	VehicleID string `json:"vehicle_id"`
	Message   string `json:"message"`
	Timestamp time.Time
}

var alerts []Alert

func main() {
	r := gin.Default()

	r.POST("/alert", func(c *gin.Context) {
		var a Alert
		if err := c.BindJSON(&a); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		a.Timestamp = time.Now()
		alerts = append(alerts, a)
		log.Printf("ALERT: %s - %s\n", a.VehicleID, a.Message)
		c.JSON(200, gin.H{"status": "sent"})
	})

	r.GET("/alerts", func(c *gin.Context) {
		c.JSON(200, alerts)
	})

	log.Println("Notification Service running on :8082")
	r.Run(":8082")
}
