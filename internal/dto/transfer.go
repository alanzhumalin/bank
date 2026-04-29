package dto

import (
	"github.com/alanzhumalin/bank/internal/domain"
)

type CreateTransferRequest struct {
	SenderAccountId   int           `json:"sender_account_id"`
	ReceiverAccountId int           `json:"receiver_account_id"`
	CurrencyId        int           `json:"currency_id"`
	Amount            int           `json:"amount"`
	Status            domain.Status `json:"status"`
}
