package redis

import (
	"fmt"
	"log/slog"
	"strconv"

	services "github.com/EvansTrein/RESTful_exchangerServer/internal/services/wallet"
	"github.com/go-redis/redis"
)

func (r *RedisDB) SetExchange(fromCurrency, toCurrency string, value float32) error {
	op := "Redis: saving the exchange rate in the cache"
	log := r.log.With(slog.String("operation", op))
	log.Debug("SetExchange func call", "fromCurrency", fromCurrency, "toCurrency", toCurrency, "value", value)

	key := fmt.Sprintf("%s/%s", fromCurrency, toCurrency)

	err := r.client.Set(key, value, r.ttlKeys).Err()
	if err != nil {
		r.log.Error("failed to save string to Redis", "error", err)
		return err
	}

	r.log.Info("exchange rate has been successfully cached", "key", key, "value", value)
	return nil
}

func (r *RedisDB) GetExchange(fromCurrency, toCurrency string) (float32, error) {
	op := "Redis: getting exchange rate from cache"
	log := r.log.With(slog.String("operation", op))
	log.Debug("GetExchange func call", "fromCurrency", fromCurrency, "toCurrency", toCurrency)

	key := fmt.Sprintf("%s/%s", fromCurrency, toCurrency)

	value, err := r.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, services.ErrRateInCacheNotFound
		}
		r.log.Error("failed to get key from Redis", "key", key, "error", err)
		return 0, err
	}

	log.Debug("cached data was retrieved", "value", value)

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		r.log.Error("failed to convert string to float", "value", value, "error", err)
		return 0, err
	}

	log.Info("the exchange rate was successfully retrieved from the cache", "key", key, "vaule", floatValue)
	return float32(floatValue), nil
}
