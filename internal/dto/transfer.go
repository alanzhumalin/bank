package dto

import (
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
)

type CreateTransferRequest struct {
	SenderAccountId   int `json:"sender_account_id"`
	ReceiverAccountId int `json:"receiver_account_id"`
	CurrencyId        int `json:"currency_id"`
	Amount            int `json:"amount"`
}

func (c *CreateTransferRequest) Validate() error {
	if c.SenderAccountId <= 0 {
		return ErrorSenderIdRequired
	}

	if c.ReceiverAccountId <= 0 {
		return ErrorReceiverIdRequired
	}

	if c.CurrencyId <= 0 {
		return ErrorCurrencyIdRequired
	}

	if c.Amount <= 0 {
		return ErrorAmountRequired
	}

	return nil
}

type TransferResponse struct {
	Id                int       `json:"id"`
	SenderAccountId   int       `json:"sender_account_id"`
	SenderFirstName   string    `json:"sender_firstname"`
	SenderLastName    string    `json:"sender_lastname"`
	ReceiverAccountId int       `json:"receiver_account_id"`
	ReceiverFirstName string    `json:"receiver_firstname"`
	ReceiverLastName  string    `json:"receiver_lastname"`
	CurrencyId        int       `json:"currency_id"`
	CurrencyName      string    `json:"currency_name"`
	CurrencyCode      string    `json:"currency_code"`
	CurrencySymbol    string    `json:"currency_symbol"`
	Amount            float64   `json:"amount"`
	Status            string    `json:"status"`
	StatusMessage     string    `json:"status_message"`
	CreatedAt         time.Time `json:"created_at"`
}

func ToTransferResponse(tr domain.Transfer) *TransferResponse {
	return &TransferResponse{
		Id:                tr.Id,
		SenderAccountId:   tr.SenderAccountId,
		SenderFirstName:   tr.Sender.FirstName,
		SenderLastName:    tr.Sender.LastName,
		ReceiverAccountId: tr.ReceiverAccountId,
		ReceiverFirstName: tr.Receiver.FirstName,
		ReceiverLastName:  tr.Receiver.LastName,
		CurrencyId:        tr.CurrencyId,
		CurrencyName:      tr.Currency.Name,
		CurrencyCode:      tr.Currency.Code,
		CurrencySymbol:    tr.Currency.Symbol,
		Amount:            tr.Amount,
		Status:            tr.Status,
		StatusMessage:     tr.StatusMessage,
		CreatedAt:         tr.CreatedAt,
	}
}
