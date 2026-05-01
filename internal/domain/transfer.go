package domain

import "github.com/shopspring/decimal"

type Transfer struct {
	Id                int
	TransactionId     int
	SenderAccountId   int
	ReceiverAccountId int
	CurrencyId        int
	Amount            decimal.Decimal
	Sender            User
	Receiver          User
	Currency          Сurrency
}
