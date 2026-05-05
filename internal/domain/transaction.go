package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Id             int
	Type           string
	Amount         decimal.Decimal
	AccountId      int
	Status         string
	StatusMessage  string
	CurrencyId     int
	CurrencyCode   string
	CurrencySymbol string

	WithDrawalDetail *WithdrawalDetail
	TransferDetail   *TransferDetail
	DepositDetail    *DepositDetail

	CreatedAt time.Time
}

type DepositDetail struct {
	Source string
}

type TransferDetail struct {
	Sender   UserDetail
	Receiver UserDetail
}

type UserDetail struct {
	FirstName   string
	LastName    string
	PhoneNumber string
}

type WithdrawalDetail struct {
	Source string
}
