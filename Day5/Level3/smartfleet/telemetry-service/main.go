
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/redis/go-redis/v9"
)

// TelemetryPayload is the JSON schema accepted by this service.
type TelemetryPayload struct {
	VehicleID string  `json:"vehicle_id"`
	Speed     float64 `json:"speed"`      // km/h
	Fuel      float64 `json:"fuel_level"` // percentage 0-100
	Lat       float64 `json:"latitude"`
	Lon       float64 `json:"longitude"`
	Ts        int64   `json:"ts"` // unix millis (optional - server will set if empty)
}

// configuration via env (12-factor)
var (
	httpAddr   = getenv("HTTP_ADDR", ":8081")
	kafkaAddr  = getenv("KAFKA_BROKER", "kafka:9092")
	kafkaTopic = getenv("KAFKA_TOPIC", "telemetry.events")
	redisAddr  = getenv("REDIS_ADDR", "redis:6379")
	redisTTL   = getenvInt("REDIS_TTL_SECONDS", 3600) // seconds
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if s := os.Getenv(key); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			return n
		}
	}
	return def
}

// global counters (atomic)
var (
	recvCounter uint64
	okCounter   uint64
	errCounter  uint64
)

// validate basic payload fields
func (t *TelemetryPayload) Validate() error {
	if t.VehicleID == "" {
		return errors.New("vehicle_id required")
	}
	// sanity checks (domain-specific)
	if t.Speed < 0 || t.Speed > 400 {
		return fmt.Errorf("speed out of range: %.2f", t.Speed)
	}
	if t.Fuel < -1 || t.Fuel > 200 { // allow -1 if not applicable
		return fmt.Errorf("fuel out of range: %.2f", t.Fuel)
	}
	// lat lon basic range
	if t.Lat < -90 || t.Lat > 90 || t.Lon < -180 || t.Lon > 180 {
		return fmt.Errorf("invalid lat/lon")
	}
	return nil
}

func main() {
	logger := log.New(os.Stdout, "[telemetry] ", log.LstdFlags|log.Lmsgprefix)

	// Kafka writer (producer)
	kWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddr},
		Topic:    kafkaTopic,
		Balancer: &kafka.Hash{}, // Keyed by vehicle id for ordering
		Async:    false,
	})
	defer func() {
		_ = kWriter.Close()
	}()

	// Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer func() {
		_ = rdb.Close()
	}()

	// Verify connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := kWriter.Close(); err != nil {
		// we created writer then closed immediately to verify; recreate real writer
		// re-create writer for normal operation
		kWriter = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{kafkaAddr},
			Topic:    kafkaTopic,
			Balancer: &kafka.Hash{},
		})
	}
	// Ping Redis
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Fatalf("redis ping failed: %v", err)
	}

	// HTTP handlers
	mux := http.NewServeMux()

	mux.HandleFunc("/telemetry", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&recvCounter, 1)
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			atomic.AddUint64(&errCounter, 1)
			return
		}
		var tp TelemetryPayload
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&tp); err != nil {
			http.Error(w, "invalid payload: "+err.Error(), http.StatusBadRequest)
			atomic.AddUint64(&errCounter, 1)
			return
		}
		// set timestamp server-side if missing or unreasonable
		if tp.Ts <= 0 {
			tp.Ts = time.Now().UnixMilli()
		}
		if err := tp.Validate(); err != nil {
			http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
			atomic.AddUint64(&errCounter, 1)
			return
		}

		// marshal payload for kafka and redis
		value, err := json.Marshal(tp)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			atomic.AddUint64(&errCounter, 1)
			return
		}

		// Write to Kafka with short timeout
		kctx, kcancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer kcancel()
		msg := kafka.Message{
			Key:   []byte(tp.VehicleID),
			Value: value,
			Time:  time.UnixMilli(tp.Ts),
		}
		if err := kWriter.WriteMessages(kctx, msg); err != nil {
			logger.Printf("kafka write failed: %v", err)
			http.Error(w, "enqueue failed", http.StatusInternalServerError)
			atomic.AddUint64(&errCounter, 1)
			return
		}

		// Update Redis latest state (non-blocking but wait for the SET result)
		redisKey := "vehicle:latest:" + tp.VehicleID
		rctx, rcancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer rcancel()
		if err := rdb.Set(rctx, redisKey, value, time.Duration(redisTTL)*time.Second).Err(); err != nil {
			// log but do not fail the request (best-effort cache)
			logger.Printf("redis set failed for %s: %v", tp.VehicleID, err)
		}

		w.WriteHeader(http.StatusAccepted)
		atomic.AddUint64(&okCounter, 1)
	})

	// health & metrics endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// simple health: redis ok, kafka writer exists
		// attempt quick redis ping
		ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancel()
		health := map[string]interface{}{
			"status":          "ok",
			"received_total":  atomic.LoadUint64(&recvCounter),
			"accepted_total":  atomic.LoadUint64(&okCounter),
			"errors_total":    atomic.LoadUint64(&errCounter),
			"redis_connected": false,
			"kafka_topic":     kafkaTopic,
		}
		if err := rdb.Ping(ctx).Err(); err == nil {
			health["redis_connected"] = true
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(health)
	})

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	// graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		// listen for OS signals
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		logger.Println("shutdown signal received: shutting down HTTP server...")

		// stop accepting requests
		ctxShut, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		if err := server.Shutdown(ctxShut); err != nil {
			logger.Printf("HTTP server Shutdown: %v", err)
		}

		// flush kafka writer (close will flush)
		logger.Println("closing kafka writer...")
		if err := kWriter.Close(); err != nil {
			logger.Printf("kafka writer close err: %v", err)
		}
		// close redis client
		if err := rdb.Close(); err != nil {
			logger.Printf("redis close err: %v", err)
		}

		close(idleConnsClosed)
	}()

	logger.Printf("starting HTTP server on %s (kafka=%s redis=%s)\n", httpAddr, kafkaAddr, redisAddr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// fatal error
		logger.Fatalf("ListenAndServe(): %v", err)
	}

	<-idleConnsClosed
	logger.Println("service stopped")
}


