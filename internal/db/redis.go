package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedisClient(addr string, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctxWithTimeout).Err(); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil

}
