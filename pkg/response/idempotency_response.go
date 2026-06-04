package response

import (
	"encoding/json"

	"github.com/alanzhumalin/bank/internal/dto"
)

func IdemErrorResponse(err error) dto.IdempotencyResponse {
	b, _ := json.Marshal(map[string]any{
		"data": err.Error(),
	})

	return dto.IdempotencyResponse{
		Status:   "failed",
		Response: json.RawMessage(b),
	}
}
