package dto

import "github.com/shopspring/decimal"

type CreateDepositRequest struct {
	Amount decimal.Decimal `json:"amount"`
	Source string          `json:"source"`
}
