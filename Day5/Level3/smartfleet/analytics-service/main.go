
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	kafkaBroker = getenv("KAFKA_BROKER", "kafka:9092")
	kafkaTopic  = getenv("KAFKA_TOPIC", "telemetry.events")
	kafkaGroup  = getenv("KAFKA_GROUP", "analytics-group")

	postgresDSN = getenv("POSTGRES_DSN", "postgres://fleet:fleetpass@postgres:5432/fleetdb?sslmode=disable")
)

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	logger := log.New(os.Stdout, "[analytics] ", log.LstdFlags|log.Lmsgprefix)

	// Setup GORM + Postgres
	db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
	if err != nil {
		logger.Fatalf("failed connect postgres: %v", err)
	}

	// connection pool tuning via sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("failed get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := MigrateSchemas(db); err != nil {
		logger.Fatalf("migrate schemas: %v", err)
	}

	store := NewStore(db, logger)

	// Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		GroupID: kafkaGroup,
		Topic:   kafkaTopic,
		MinBytes: 1e3,   // 1KB
		MaxBytes: 10e6,  // 10MB
	})
	defer reader.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// run consumer loop
	logger.Println("starting consumer loop...")
	if err := runConsumerLoop(ctx, reader, store, logger); err != nil {
		logger.Fatalf("consumer loop ended with error: %v", err)
	}

	logger.Println("analytics service stopped gracefully")
}


