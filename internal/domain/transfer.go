package domain

import "time"

type Status string

var (
	Pending   Status = "pending"
	Completed Status = "completed"
	Failed    Status = "failed"
)

type Transfer struct {
	Id                int
	SenderAccountId   int
	ReceiverAccountId int
	CurrencyId        int
	Amount            int
	Status            Status
	CreatedAt         time.Time
}

func NewTransferFromDb(id int, senderId int, receiverId int, currencyId int, amount int, status Status, created_at time.Time) Transfer {
	return Transfer{
		Id:                id,
		SenderAccountId:   senderId,
		ReceiverAccountId: receiverId,
		CurrencyId:        currencyId,
		Amount:            amount,
		Status:            status,
		CreatedAt:         created_at,
	}
}

func NewTransfer(id int, senderId int, receiverId int, currencyId int, amount int, status Status, created_at time.Time) Transfer {
	return Transfer{
		SenderAccountId:   senderId,
		ReceiverAccountId: receiverId,
		CurrencyId:        currencyId,
		Amount:            amount,
		Status:            status,
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
