
package kafka

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// Configurable via env
var (
	KafkaBroker = getenv("KAFKA_BROKER", "localhost:9092")
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// NewProducer returns a Kafka writer (producer)
func NewProducer(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(KafkaBroker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}
}

// ProduceMessage publishes a message to Kafka
func ProduceMessage(ctx context.Context, writer *kafka.Writer, key []byte, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return writer.WriteMessages(ctx, msg)
}

// NewConsumer returns a Kafka reader (consumer)
func NewConsumer(topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{KafkaBroker},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

// ConsumeMessages reads messages and invokes handler
func ConsumeMessages(ctx context.Context, reader *kafka.Reader, handler func(kafka.Message) error) {
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("[kafka] fetch error: %v", err)
			continue
		}
		if err := handler(msg); err != nil {
			log.Printf("[kafka] handler error: %v", err)
		}
		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("[kafka] commit error: %v", err)
		}
	}
}


