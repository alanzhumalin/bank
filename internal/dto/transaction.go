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

type DepositDetail struct {
	Source string `json:"source"`
}

type TransferDetail struct {
	Sender   UserDetail `json:"sender"`
	Receiver UserDetail `json:"receiver"`
}

type UserDetail struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	PhoneNumber string `json:"phone_number"`
}

type WithdrawalDetail struct {
	Source string `json:"source"`
}

func ToTransactionResponse(tr domain.Transaction) TransactionResponse {

	var detail any

	if tr.WithDrawalDetail != nil {
		detail = WithdrawalDetail{
			Source: tr.WithDrawalDetail.Source,
		}
	}

	if tr.DepositDetail != nil {
		detail = DepositDetail{
			Source: tr.DepositDetail.Source,
		}
	}

	if tr.TransferDetail != nil {
		detail = TransferDetail{
			Sender: UserDetail{
				FirstName:   tr.TransferDetail.Sender.FirstName,
				LastName:    tr.TransferDetail.Sender.LastName,
				PhoneNumber: tr.TransferDetail.Sender.PhoneNumber,
			},
			Receiver: UserDetail{
				FirstName:   tr.TransferDetail.Receiver.FirstName,
				LastName:    tr.TransferDetail.Receiver.LastName,
				PhoneNumber: tr.TransferDetail.Receiver.PhoneNumber,
			},
		}
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
