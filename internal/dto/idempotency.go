package dto

import (
	"encoding/json"

	"github.com/alanzhumalin/bank/internal/domain"
)

type IdempotencyResponse struct {
	Status   string          `json:"status"`
	Response json.RawMessage `json:"response"`
}

func ToIdempotency(key string, userId int, operation string) domain.Idempotency {
	return domain.Idempotency{
		IdempotencyKey: key,
		UserId:         userId,
		Operation:      operation,
	}
}
