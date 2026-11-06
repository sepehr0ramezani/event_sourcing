package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectToRedis(ctx context.Context) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
		Protocol: 3,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}
