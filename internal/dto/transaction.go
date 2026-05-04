package dto

import (
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/shopspring/decimal"
)

type TransactionResponse struct {
	Id            int             `json:"id"`
	Type          string          `json:"type"`
	Amount        decimal.Decimal `json:"amount"`
	AccountId     int             `json:"account_id"`
	Status        string          `json:"status"`
	StatusMessage string          `json:"status_message"`
	CreatedAt     time.Time       `json:"created_at"`
}

func ToTransactionResponse(tr domain.Transaction) TransactionResponse {
	return TransactionResponse{
		Id:            tr.Id,
		Type:          tr.Type,
		Amount:        tr.Amount,
		AccountId:     tr.AccountId,
		Status:        tr.Status,
		StatusMessage: tr.StatusMessage,
		CreatedAt:     tr.CreatedAt,
	}
}
