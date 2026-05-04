package dto

import (
	"time"

	"github.com/alanzhumalin/bank/internal/domain"
	"github.com/shopspring/decimal"
)

type CreateTransferRequest struct {
	SenderAccountId   int             `json:"sender_account_id"`
	ReceiverAccountId int             `json:"receiver_account_id"`
	CurrencyId        int             `json:"currency_id"`
	Amount            decimal.Decimal `json:"amount"`
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

	if c.Amount.LessThanOrEqual(decimal.Zero) {
		return ErrorAmountRequired
	}

	return nil
}

func NewTransferFromDb(id int, senderId int, receiverId int, currency domain.Сurrency, currencyId int, amount decimal.Decimal, sender domain.User, receiver domain.User, status string, created_at time.Time, statusMessage string) domain.Transfer {
	return domain.Transfer{
		Id:                id,
		SenderAccountId:   senderId,
		ReceiverAccountId: receiverId,
		CurrencyId:        currencyId,
		Amount:            amount,
		Sender:            sender,
		Receiver:          receiver,
		Currency:          currency,
	}
}

func NewTransfer(senderId int, receiverId int, currencyId int, amount decimal.Decimal) domain.Transfer {
	return domain.Transfer{
		SenderAccountId:   senderId,
		ReceiverAccountId: receiverId,
		CurrencyId:        currencyId,
		Amount:            amount,
	}
}

/*
CREATE type transfer_status as enum ('pending', 'completed', 'failed');

CREATE TABLE IF NOT EXISTS transfers(
    id BIGINT generated always as identity PRIMARY key,
    sender_account_id BIGINT not null REFERENCES accounts(id),
    receiver_account_id BIGINT not null REFERENCES accounts(id),
    currency_id BIGINT not null REFERENCES currencies(id),
    amount numeric(12,2) not null check (amount > 0),
    status transfer_status not null DEFAULT 'pending',
    created_at TIMESTAMPtz not null default now()
);
*/
