package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Vehicle struct {
	ID          string  `json:"id"`
	Speed       float64 `json:"speed"`
	Temperature float64 `json:"temperature"`
	Status      string  `json:"status"`
}

var vehicles = []Vehicle{

	{"VW001", 72.4, 36.8, "Running"},
	{"VW002", 0.0, 28.1, "Idle"},
}

func updateVehicleData() {

	for {
		for i := range vehicles {
			vehicles[i].Speed = rand.Float64() * 120
			vehicles[i].Temperature = 25 + rand.Float64()*20
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {

	go updateVehicleData()

	router := gin.Default()

	router.GET("/vehicles", func(c *gin.Context) {
		c.JSON(http.StatusOK, vehicles)
	})

	router.GET("/vehicle/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, v := range vehicles {
			if v.ID == id {
				c.JSON(http.StatusOK, v)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
	})

	router.Run(":8082")
}
