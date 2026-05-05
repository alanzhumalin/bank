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
	Details       any             `json:"details"`
	CreatedAt     time.Time       `json:"created_at"`
}

func ToTransactionResponse(tr domain.Transaction) TransactionResponse {

	var detail any

	if tr.WithDrawalDetail != nil {
		detail = tr.WithDrawalDetail
	}

	if tr.DepositDetail != nil {
		detail = tr.WithDrawalDetail
	}

	if tr.TransferDetail != nil {
		detail = tr.TransferDetail
	}

	return TransactionResponse{
		Id:            tr.Id,
		Type:          tr.Type,
		Amount:        tr.Amount,
		AccountId:     tr.AccountId,
		Status:        tr.Status,
		StatusMessage: tr.StatusMessage,
		Details:       detail,
		CreatedAt:     tr.CreatedAt,
	}
}
