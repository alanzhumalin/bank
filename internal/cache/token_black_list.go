package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var JTIIsRequired = errors.New("Jti is required")

type tokenBlackList struct {
	redisClient *redis.Client
}

func NewTokenBlackList(redisClient *redis.Client) (TokenBlackList, error) {
	if redisClient == nil {
		return nil, RedisClientIsRequired
	}

	return &tokenBlackList{
		redisClient: redisClient,
	}, nil
}

func (t *tokenBlackList) Add(ctx context.Context, jti string, ttl time.Duration) error {
	if jti == "" {
		return JTIIsRequired
	}

	if ttl <= 0 {
		return nil
	}
	key := fmt.Sprintf("blacklist:access:%s", jti)
	if err := t.redisClient.Set(ctx, key, "access_token", ttl).Err(); err != nil {
		return fmt.Errorf("Error in adding key to redis: %w", err)
	}

	return nil
}

func (t *tokenBlackList) Exists(ctx context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, JTIIsRequired
	}
	key := fmt.Sprintf("blacklist:access:%s", jti)

	n, err := t.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("Error in checking existence of the key: %w", err)
	}

	return n > 0, nil
}
