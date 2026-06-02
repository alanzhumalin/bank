package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisErr = errors.New("Error occured with redis")

type IdempotencyStore struct {
	redisClient   *redis.Client
	ttlProcessing time.Duration
	ttlCompleted  time.Duration
}

func NewIdempotencyStore(redisClient *redis.Client, ttlProcessing time.Duration, ttlCompleted time.Duration) (*IdempotencyStore, error) {
	if redisClient == nil {
		return nil, RedisClientIsRequired
	}

	if ttlCompleted <= 0 {
		return nil, TTLIsRequired
	}

	if ttlProcessing <= 0 {
		return nil, TTLIsRequired
	}

	return &IdempotencyStore{
		redisClient:   redisClient,
		ttlProcessing: ttlProcessing,
		ttlCompleted:  ttlCompleted,
	}, nil
}

func (iS *IdempotencyStore) Start(ctx context.Context, key string, status string) (bool, string, error) {
	ok, err := iS.redisClient.SetNX(ctx, key, status, iS.ttlProcessing).Result()

	if err != nil {
		return false, "", fmt.Errorf("Error occured while setting the key in redis: %w", err)
	}

	if !ok {
		res, err := iS.redisClient.Get(ctx, key).Result()
		if err != nil {
			return false, "", fmt.Errorf("Error occured while getting the result in redis: %w", err)
		}
		return false, res, nil
	}

	return true, status, nil
}

func (iS *IdempotencyStore) Complete(ctx context.Context, key string, response string) error {
	if err := iS.redisClient.Set(ctx, key, response, iS.ttlCompleted).Err(); err != nil {
		return fmt.Errorf("Error in completing the status of the key: %w", err)
	}
	return nil
}
