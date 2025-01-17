package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

func (r *RedisDB) TestMethodSet(key, value string) error {

	err := r.client.Set(key, value, r.ttlKeys).Err()
	if err != nil {
		r.log.Error("failed to save string to Redis", "key", key, "error", err)
		return fmt.Errorf("failed to save string to Redis: %w", err)
	}

	r.log.Warn("string saved to Redis", "key", key)
	return nil
}

func (r *RedisDB) TestMethodGet(key string) (string, error) {

	val, err := r.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			r.log.Warn("key not found in Redis", "key", key)
			return "", fmt.Errorf("key not found: %s", key)
		}
		r.log.Error("failed to get string from Redis", "key", key, "error", err)
		return "", fmt.Errorf("failed to get string from Redis: %w", err)
	}

	r.log.Warn("string retrieved from Redis", "key", key, "value", val)
	return val, nil
}
