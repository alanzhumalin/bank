package domain

import "github.com/shopspring/decimal"

type Deposit struct {
	Id            int
	TransactionId int
	AccountId     int
	Amount        decimal.Decimal
	Source        string
}
