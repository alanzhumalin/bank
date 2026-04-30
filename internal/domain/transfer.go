package domain

import (
	"time"
)

type Transfer struct {
	Id                int
	SenderAccountId   int
	ReceiverAccountId int
	CurrencyId        int
	Amount            float64
	Status            string
	CreatedAt         time.Time
	StatusMessage     string
	Sender            User
	Receiver          User
	Currency          Сurrency
}

func NewTransferFromDb(id int, senderId int, receiverId int, currency Сurrency, currencyId int, amount float64, sender User, receiver User, status string, created_at time.Time, statusMessage string) Transfer {
	return Transfer{
		Id:                id,
		SenderAccountId:   senderId,
		ReceiverAccountId: receiverId,
		CurrencyId:        currencyId,
		Amount:            amount,
		Status:            status,
		CreatedAt:         created_at,
		StatusMessage:     statusMessage,
		Sender:            sender,
		Receiver:          receiver,
		Currency:          currency,
	}
}

func NewTransfer(senderId int, receiverId int, currencyId int, amount float64) Transfer {
	return Transfer{
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
