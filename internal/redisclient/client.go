// Package redisclient wraps the go-redis client for use in stmt2redis.
package redisclient

import (
	"context"
	"fmt"
	"os"

	"github.com/its-the-vibe/stmt2redis/internal/config"
	"github.com/redis/go-redis/v9"
)

// Client wraps a go-redis client.
type Client struct {
	rdb *redis.Client
}

// New creates a new Redis client using the supplied configuration. The Redis
// password is read from the REDIS_PASSWORD environment variable (loaded from
// .env by the caller before this function is invoked).
func New(cfg *config.Config) (*Client, error) {
	password := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connecting to Redis at %s:%d: %w", cfg.Redis.Host, cfg.Redis.Port, err)
	}

	return &Client{rdb: rdb}, nil
}

// RPush appends values to the tail of the Redis list identified by key.
func (c *Client) RPush(ctx context.Context, key string, values ...string) error {
	args := make([]interface{}, len(values))
	for i, v := range values {
		args[i] = v
	}
	if err := c.rdb.RPush(ctx, key, args...).Err(); err != nil {
		return fmt.Errorf("RPUSH to list %q: %w", key, err)
	}
	return nil
}

// Close closes the underlying Redis connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}
