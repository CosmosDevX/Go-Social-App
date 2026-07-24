package infrastructure

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() RedisClient {
	return RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:        "localhost:6379",
			Password:    "",
			DB:          0,
			PoolSize:    10,
			PoolTimeout: 30 * time.Second,
		}),
	}
}

func (c RedisClient) Setup(ctx context.Context) {
	_, err := c.client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
}

func (c RedisClient) GetClient() *redis.Client {
	return c.client
}

func (c RedisClient) Shutdown() error {
	if err := c.client.Close(); err != nil {
		return err
	}

	return nil
}
