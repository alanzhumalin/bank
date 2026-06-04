package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/redis/go-redis/v9"
)

var RedisErr = errors.New("Error occured with redis")

type IdempotencyStore struct {
	redisClient   *redis.Client
	ttlProcessing time.Duration
	ttlFailed     time.Duration
	ttlCompleted  time.Duration
}

func NewIdempotencyStore(redisClient *redis.Client, ttlFailed time.Duration, ttlProcessing time.Duration, ttlCompleted time.Duration) (IdempotencyStoreInterface, error) {
	if redisClient == nil {
		return nil, RedisClientIsRequired
	}

	if ttlCompleted <= 0 {
		return nil, TTLIsRequired
	}

	if ttlProcessing <= 0 {
		return nil, TTLIsRequired
	}

	if ttlFailed <= 0 {
		return nil, TTLIsRequired
	}

	return &IdempotencyStore{
		redisClient:   redisClient,
		ttlProcessing: ttlProcessing,
		ttlCompleted:  ttlCompleted,
		ttlFailed:     ttlFailed,
	}, nil
}

func (iS *IdempotencyStore) Start(ctx context.Context, key string, res dto.IdempotencyResponse) (bool, dto.IdempotencyResponse, error) {
	b, err := json.Marshal(res)
	if err != nil {
		return false, dto.IdempotencyResponse{}, fmt.Errorf("Error in converting to json: %w", err)
	}

	ok, err := iS.redisClient.SetNX(ctx, key, b, iS.ttlProcessing).Result()

	if err != nil {
		return false, dto.IdempotencyResponse{}, fmt.Errorf("Error occured while setting the key in redis: %w", err)
	}

	if !ok {
		var i dto.IdempotencyResponse
		res, err := iS.redisClient.Get(ctx, key).Bytes()
		if err != nil {
			return false, dto.IdempotencyResponse{}, fmt.Errorf("Error occured while getting the result in redis: %w", err)
		}

		if err := json.Unmarshal(res, &i); err != nil {
			return false, dto.IdempotencyResponse{}, fmt.Errorf("Error in json convertation: %w", err)
		}

		return false, i, nil
	}

	return true, res, nil
}

func (iS *IdempotencyStore) Complete(ctx context.Context, key string, res dto.IdempotencyResponse) error {
	b, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("Error in converting to json: %w", err)
	}
	if err := iS.redisClient.Set(ctx, key, b, iS.ttlCompleted).Err(); err != nil {
		return fmt.Errorf("Error in completing the status of the key: %w", err)
	}
	return nil
}

func (iS *IdempotencyStore) Failed(ctx context.Context, key string, res dto.IdempotencyResponse) error {
	b, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("Error in converting to json: %w", err)
	}
	if err := iS.redisClient.Set(ctx, key, b, iS.ttlFailed).Err(); err != nil {
		return fmt.Errorf("Error occured while setting the value: %w", err)
	}

	return nil
}

func (iS *IdempotencyStore) Delete(ctx context.Context, key string) error {
	if err := iS.redisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("Error in delete key from redis: %w", err)
	}

	return nil
}
