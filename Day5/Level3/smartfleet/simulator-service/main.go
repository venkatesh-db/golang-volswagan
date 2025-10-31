
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

// VehicleTelemetry represents one vehicle's telemetry data
type VehicleTelemetry struct {
	VehicleID   string  `json:"vehicle_id"`
	SpeedKmph   float64 `json:"speed_kmph"`
	FuelPercent float64 `json:"fuel_percent"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	EngineTemp  float64 `json:"engine_temp"`
	Timestamp   int64   `json:"timestamp"`
}

// Configurable parameters
var (
	telemetryURL  = getenv("TELEMETRY_URL", "http://localhost:8081/telemetry")
	vehicleCount  = getenvInt("VEHICLE_COUNT", 5)
	sendInterval  = time.Duration(getenvInt("SEND_INTERVAL_MS", 2000)) * time.Millisecond
	speedMin      = 0.0
	speedMax      = 180.0
	fuelMin       = 5.0
	fuelMax       = 100.0
	latMin        = 12.90
	latMax        = 13.10
	lngMin        = 77.55
	lngMax        = 77.70
	engineMinTemp = 70.0
	engineMaxTemp = 120.0
)

// helper functions
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// generate random float between min and max
func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// simulate one vehicle's telemetry
func simulateVehicle(ctx context.Context, vehicleID string, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(sendInterval)
	defer ticker.Stop()

	client := &http.Client{Timeout: 2 * time.Second}

	for {
		select {
		case <-ctx.Done():
			log.Printf("[sim] vehicle %s stopped", vehicleID)
			return
		case <-ticker.C:
			data := VehicleTelemetry{
				VehicleID:   vehicleID,
				SpeedKmph:   randomFloat(speedMin, speedMax),
				FuelPercent: randomFloat(fuelMin, fuelMax),
				Latitude:    randomFloat(latMin, latMax),
				Longitude:   randomFloat(lngMin, lngMax),
				EngineTemp:  randomFloat(engineMinTemp, engineMaxTemp),
				Timestamp:   time.Now().Unix(),
			}
			payload, _ := json.Marshal(data)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, telemetryURL, bytes.NewReader(payload))
			if err != nil {
				log.Printf("[sim] vehicle %s error creating request: %v", vehicleID, err)
				continue
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("[sim] vehicle %s error sending telemetry: %v", vehicleID, err)
				continue
			}
			resp.Body.Close()
			log.Printf("[sim] vehicle %s sent telemetry: speed=%.1f kmph fuel=%.1f%% lat=%.5f lng=%.5f temp=%.1f°C",
				vehicleID, data.SpeedKmph, data.FuelPercent, data.Latitude, data.Longitude, data.EngineTemp)
		}
	}
}

func main() {
	log.Println("[sim] Simulator Service starting...")
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// start multiple vehicle goroutines
	for i := 0; i < vehicleCount; i++ {
		vehicleID := uuid.New().String()
		wg.Add(1)
		go simulateVehicle(ctx, vehicleID, wg)
	}

	// graceful shutdown on signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("[sim] shutdown signal received...")
	cancel()
	wg.Wait()
	log.Println("[sim] Simulator Service stopped")
}


