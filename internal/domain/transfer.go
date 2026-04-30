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
