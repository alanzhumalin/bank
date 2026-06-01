package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClientIsRequired = errors.New("Redis client is required")
	LimitIsRequired       = errors.New("Limit is required")
	TTLIsRequired         = errors.New("Ttl is required")
)

type RateLimiter struct {
	ClientRedis *redis.Client
	Limit       int64
	TTL         time.Duration
}

func NewRateLimiter(client *redis.Client, limit int64, ttl time.Duration) (*RateLimiter, error) {

	if client == nil {
		return nil, RedisClientIsRequired
	}

	if limit <= 0 {
		return nil, LimitIsRequired
	}

	if ttl < 0 {
		return nil, TTLIsRequired
	}

	return &RateLimiter{ClientRedis: client, Limit: limit, TTL: ttl}, nil
}

func (r *RateLimiter) IsAllowed(ctx context.Context, key string) (bool, error) {

	count, err := r.ClientRedis.Incr(ctx, key).Result()

	if err != nil {
		return false, fmt.Errorf("Error in isallowed: %w", err)
	}

	if count == 1 {
		if err := r.ClientRedis.Expire(ctx, key, r.TTL).Err(); err != nil {
			return false, err
		}
	}

	return count <= r.Limit, nil
}
