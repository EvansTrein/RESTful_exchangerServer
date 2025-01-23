package redis

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-redis/redis"
)

// RedisDB represents a connection to a Redis database.
// It includes a Redis client, a logger, and a TTL (Time-To-Live) for keys.
type RedisDB struct {
	client  *redis.Client
	log     *slog.Logger
	ttlKeys time.Duration
}

// New creates a new RedisDB instance and establishes a connection to the Redis server.
// It takes the host, port, password, and TTL for keys as parameters.
// If the connection fails, it returns an error.
func New(log *slog.Logger, host string, port string, password string, ttlKeys time.Duration) (*RedisDB, error) {
	log.Debug("redis: connection to Redis started")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Info("redis: connect to Redis successfully")
	return &RedisDB{client: client, ttlKeys: ttlKeys, log: log}, nil
}

// Close closes the connection to the Redis server.
// If the connection is already closed, it returns an error.
func (r *RedisDB) Close() error {
	r.log.Debug("redis: stop started")

	if r.client == nil {
		return fmt.Errorf("redis connection is already closed")
	}

	if err := r.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	r.client = nil

	r.log.Info("redis: stop successful")
	return nil
}
