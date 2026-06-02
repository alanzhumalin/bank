package dto

import "github.com/alanzhumalin/bank/internal/domain"

func ToIdempotency(key string, userId int, operation string) domain.Idempotency {
	return domain.Idempotency{
		IdempotencyKey: key,
		UserId:         userId,
		Operation:      operation,
	}
}
