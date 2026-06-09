package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/alanzhumalin/bank/internal/dto"
	"github.com/redis/go-redis/v9"
)

type otpStore struct {
	redisClient *redis.Client
	ttlDuration time.Duration
}

func NewOTPStore(redisClient *redis.Client, ttlDuration time.Duration) (*otpStore, error) {
	if redisClient == nil {
		return nil, RedisClientIsRequired
	}

	if ttlDuration <= 0 {
		return nil, TTLIsRequired
	}

	return &otpStore{
		redisClient: redisClient,
		ttlDuration: ttlDuration,
	}, nil
}

func (o *otpStore) Save(ctx context.Context, event string, challengeId string, detail domain.OTPDetail) error {
	key := fmt.Sprintf("otp:%s:%s", event, challengeId)

	b, err := json.Marshal(detail)

	if err != nil {
		return fmt.Errorf("error in marshaling detail of otp: %w", err)
	}

	if err := o.redisClient.Set(ctx, key, b, o.ttlDuration).Err(); err != nil {
		return fmt.Errorf("error in saving key with has for otp: %w", err)
	}

	return nil
}

func (o *otpStore) Verify(ctx context.Context, event string, challengeId string, codeHash string) (bool, string, error) {

	if challengeId == "" {
		return false, "", dto.ChallengeIdIsRequired
	}

	if codeHash == "" {
		return false, "", dto.CodeHashIsRequired
	}

	key := fmt.Sprintf("otp:%s:%s", event, challengeId)

	b, err := o.redisClient.Get(ctx, key).Bytes()

	if err != nil {
		return false, "", fmt.Errorf("error in get the key from redis for otp: %w", err)
	}

	var detail domain.OTPDetail

	if err := json.Unmarshal(b, &detail); err != nil {
		return false, "", fmt.Errorf("error in unmarshaling the bytes from redis otp: %w", err)
	}

	if detail.CodeHash != codeHash {

		return false, "", nil
	}

	return true, detail.PhoneNumber, nil
}

func (o *otpStore) Delete(ctx context.Context, key string) error {
	if err := o.redisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("error in deleting key from redis: %w", err)
	}

	return nil
}
