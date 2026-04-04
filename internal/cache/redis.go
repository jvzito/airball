package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/jvzito/airball/internal/config"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func New(cfg *config.Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &Cache{client: client, ttl: time.Duration(cfg.CacheTTL) * time.Second}, nil
}

func (c *Cache) Get(ctx context.Context, key string, dest any) (bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal([]byte(val), dest)
}

func (c *Cache) Set(ctx context.Context, key string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, b, c.ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) Close() error { return c.client.Close() }
