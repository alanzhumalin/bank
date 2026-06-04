package cache

import (
	"context"
	"time"

	"github.com/alanzhumalin/bank/internal/dto"
)

type IdempotencyStoreInterface interface {
	Start(ctx context.Context, key string, res dto.IdempotencyResponse) (bool, dto.IdempotencyResponse, error)
	Complete(ctx context.Context, key string, res dto.IdempotencyResponse) error
	Failed(ctx context.Context, key string, res dto.IdempotencyResponse) error
	Delete(ctx context.Context, key string) error
}

type TokenBlackList interface {
	Add(ctx context.Context, jti string, ttl time.Duration) error
	Exists(ctx context.Context, jti string) (bool, error)
}
