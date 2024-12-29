package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisCache struct {
	ctx     context.Context
	client  *redis.Client
	timeout time.Duration
}

var (
	RedisHost = config.GetString("REDIS_HOST", "localhost")
	RedisPort = config.GetInt("REDIS_PORT", 6379)
)

// Option defines a function type for configuring the RedisCache.
type Option func(*RedisCache)

// WithTimeout sets the timeout for Redis operations.
func WithTimeout(timeout time.Duration) Option {
	return func(c *RedisCache) {
		c.timeout = timeout
	}
}

func NewRedisCache(options ...Option) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", RedisHost, RedisPort),
		Password: "",
		DB:       0,
	})

	// Test the connection to the Redis server
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Printf("failed to connect to Redis: %v", err)
		return nil
	}

	cache := &RedisCache{
		client:  client,
		ctx:     context.Background(),
		timeout: 5 * time.Second,
	}

	// Apply custom options
	for _, option := range options {
		option(cache)
	}

	return cache
}

// Set stores a key-value pair in the cache with an expiration duration.
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("failed to marshal cache value: %v", err)
		return err
	}

	if err = r.client.Set(ctx, key, data, expiration).Err(); err != nil {
		log.Printf("failed to set cache for key %s: %v", key, err)
		return err
	}

	return nil
}

// Get retrieves the value associated with a given key from the cache.
func (r *RedisCache) Get(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	data, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("cache miss for key %s", key)
	} else if err != nil {
		return fmt.Errorf("failed to get cache for key %s: %w", key, err)
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	return nil
}

// Delete removes a key from the cache.
func (r *RedisCache) Delete(key string) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	if _, err := r.client.Del(ctx, key).Result(); err != nil {
		return fmt.Errorf("failed to delete cache for key %s: %w", key, err)
	}

	return nil
}

// Close gracefully closes the Redis client connection.
func (r *RedisCache) Close() error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}
	return nil
}
