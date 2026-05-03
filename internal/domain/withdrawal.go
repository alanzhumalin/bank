package domain

import "github.com/shopspring/decimal"

type Withdrawal struct {
	Id            int
	TransactionId int
	AccountId     int
	Amount        decimal.Decimal
	Source        string
}
