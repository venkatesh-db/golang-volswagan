
package redis

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisAddr = getenv("REDIS_ADDR", "localhost:6379")

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// NewClient creates a Redis client with connection pool
func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         RedisAddr,
		Password:     "", // use env if needed
		DB:           0,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
}

// Ping tests Redis connectivity
func Ping(ctx context.Context, rdb *redis.Client) error {
	return rdb.Ping(ctx).Err()
}
