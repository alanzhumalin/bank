package dto

import "github.com/shopspring/decimal"

type CreateDepositRequest struct {
	Amount decimal.Decimal `json:"amount"`
	Source string          `json:"source"`
}

func (c *CreateDepositRequest) Validate() error {
	if c.Amount.IsZero() {
		return ErrorAmountIsRequired
	}

	if c.Source == "" {
		return ErrorDepositSourceIsRequired
	}

	return nil
}
